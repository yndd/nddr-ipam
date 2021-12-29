/*
Copyright 2021 NDD.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package ipamnetworkinstance

import (
	"context"
	"strings"
	"time"

	"github.com/hansthienpondt/goipam/pkg/table"
	"github.com/pkg/errors"
	nddv1 "github.com/yndd/ndd-runtime/apis/common/v1"
	"github.com/yndd/ndd-runtime/pkg/event"
	"github.com/yndd/ndd-runtime/pkg/logging"
	"github.com/yndd/ndd-runtime/pkg/meta"
	"github.com/yndd/ndd-runtime/pkg/resource"
	ipamv1alpha1 "github.com/yndd/nddr-ipam/apis/ipam/v1alpha1"
	"github.com/yndd/nddr-ipam/internal/shared"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

const (
	finalizerName = "finalizer.ipamnetworkinstance.ipam.nddr.yndd.io"
	//
	reconcileTimeout = 1 * time.Minute
	longWait         = 1 * time.Minute
	mediumWait       = 30 * time.Second
	shortWait        = 15 * time.Second
	veryShortWait    = 5 * time.Second

	// Errors
	errGetK8sResource = "cannot get ipamnetworkinstance resource"
	errUpdateStatus   = "cannot update status of ipamnetworkinstance resource"

	// events
	reasonReconcileSuccess      event.Reason = "ReconcileSuccess"
	reasonCannotDelete          event.Reason = "CannotDeleteResource"
	reasonCannotAddFInalizer    event.Reason = "CannotAddFinalizer"
	reasonCannotDeleteFInalizer event.Reason = "CannotDeleteFinalizer"
	reasonCannotInitialize      event.Reason = "CannotInitializeResource"
	reasonCannotGetAllocations  event.Reason = "CannotGetAllocations"
	reasonAppLogicFailed        event.Reason = "ApplogicFailed"
	reasonCannotGarbageCollect  event.Reason = "CannotGarbageCollect"
)

// ReconcilerOption is used to configure the Reconciler.
type ReconcilerOption func(*Reconciler)

// Reconciler reconciles packages.
type Reconciler struct {
	client  resource.ClientApplicator
	log     logging.Logger
	record  event.Recorder
	managed mrManaged

	newIpamNetworkInstance func() ipamv1alpha1.In
	newIpam                func() ipamv1alpha1.Ip
	//newIpamNetworkInstanceIpPrefixList func() ipamv1alpha1.IpPrefixList
	newAllocList func() ipamv1alpha1.AaList

	iptree map[string]*table.RouteTable
}

type mrManaged struct {
	resource.Finalizer
}

// WithLogger specifies how the Reconciler should log messages.
func WithLogger(log logging.Logger) ReconcilerOption {
	return func(r *Reconciler) {
		r.log = log
	}
}

func WithNewResourceFn(f func() ipamv1alpha1.In) ReconcilerOption {
	return func(r *Reconciler) {
		r.newIpamNetworkInstance = f
	}
}

func WithIpTree(iptree map[string]*table.RouteTable) ReconcilerOption {
	return func(r *Reconciler) {
		r.iptree = iptree
	}
}

func WithNewIpamFn(f func() ipamv1alpha1.Ip) ReconcilerOption {
	return func(r *Reconciler) {
		r.newIpam = f
	}
}

func WithNewAllocFn(f func() ipamv1alpha1.AaList) ReconcilerOption {
	return func(r *Reconciler) {
		r.newAllocList = f
	}
}

// WithRecorder specifies how the Reconciler should record Kubernetes events.
func WithRecorder(er event.Recorder) ReconcilerOption {
	return func(r *Reconciler) {
		r.record = er
	}
}

func defaultMRManaged(m ctrl.Manager) mrManaged {
	return mrManaged{
		Finalizer: resource.NewAPIFinalizer(m.GetClient(), finalizerName),
	}
}

// Setup adds a controller that reconciles ipam.
func Setup(mgr ctrl.Manager, o controller.Options, nddcopts *shared.NddControllerOptions) error {
	name := "nddr/" + strings.ToLower(ipamv1alpha1.IpamNetworkInstanceGroupKind)
	fn := func() ipamv1alpha1.In { return &ipamv1alpha1.IpamNetworkInstance{} }
	ipamfn := func() ipamv1alpha1.Ip { return &ipamv1alpha1.Ipam{} }
	afn := func() ipamv1alpha1.AaList { return &ipamv1alpha1.AllocList{} }

	r := NewReconciler(mgr,
		WithLogger(nddcopts.Logger.WithValues("controller", name)),
		WithNewResourceFn(fn),
		WithIpTree(nddcopts.Iptree),
		WithNewIpamFn(ipamfn),
		WithNewAllocFn(afn),
		WithRecorder(event.NewAPIRecorder(mgr.GetEventRecorderFor(name))),
	)

	ipamHandler := &EnqueueRequestForAllIpams{
		client: mgr.GetClient(),
		log:    nddcopts.Logger,
		ctx:    context.Background(),
		//iptree: r.iptree,
	}

	return ctrl.NewControllerManagedBy(mgr).
		Named(name).
		WithOptions(o).
		For(&ipamv1alpha1.IpamNetworkInstance{}).
		//Owns(&ipamv1alpha1.IpamNetworkInstanceIpPrefix{}).
		Watches(&source.Kind{Type: &ipamv1alpha1.Ipam{}}, ipamHandler).
		//Watches(&source.Kind{Type: &ipamv1alpha1.IpamNetworkInstanceIpPrefix{}}, ipprefixHandler).
		//Watches(&source.Kind{Type: &ipamv1alpha1.Alloc{}}, allocHandler).
		WithEventFilter(resource.IgnoreUpdateWithoutGenerationChangePredicate()).
		Complete(r)
}

// NewReconciler creates a new reconciler.
func NewReconciler(mgr ctrl.Manager, opts ...ReconcilerOption) *Reconciler {

	r := &Reconciler{
		client: resource.ClientApplicator{
			Client:     mgr.GetClient(),
			Applicator: resource.NewAPIPatchingApplicator(mgr.GetClient()),
		},
		log:     logging.NewNopLogger(),
		record:  event.NewNopRecorder(),
		managed: defaultMRManaged(mgr),
	}

	for _, f := range opts {
		f(r)
	}

	return r
}

// Reconcile ipam allocation.
func (r *Reconciler) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) { // nolint:gocyclo
	log := r.log.WithValues("request", req)
	log.Debug("Reconciling IpamNetworkInstance", "NameSpace", req.NamespacedName)

	ctx, cancel := context.WithTimeout(ctx, reconcileTimeout)
	defer cancel()

	cr := r.newIpamNetworkInstance()
	if err := r.client.Get(ctx, req.NamespacedName, cr); err != nil {
		// There's no need to requeue if we no longer exist. Otherwise we'll be
		// requeued implicitly because we return an error.
		log.Debug("Cannot get managed resource", "error", err)
		return reconcile.Result{}, errors.Wrap(resource.IgnoreNotFound(err), errGetK8sResource)
	}
	record := r.record.WithAnnotations("name", cr.GetAnnotations()[cr.GetName()])

	crName := strings.Join([]string{cr.GetNamespace(), cr.GetOrganizationName(), cr.GetIpamName(), cr.GetNetworkInstanceName()}, ".")

	if meta.WasDeleted(cr) {
		log = log.WithValues("deletion-timestamp", cr.GetDeletionTimestamp())

		// TODO check allocations

		if err := r.managed.RemoveFinalizer(ctx, cr); err != nil {
			// If this is the first time we encounter this issue we'll be
			// requeued implicitly when we update our status with the new error
			// condition. If not, we requeue explicitly, which will trigger
			// backoff.
			record.Event(cr, event.Warning(reasonCannotDeleteFInalizer, err))
			log.Debug("Cannot remove managed resource finalizer", "error", err)
			cr.SetConditions(nddv1.ReconcileError(err), ipamv1alpha1.NotReady())
			return reconcile.Result{Requeue: true}, errors.Wrap(r.client.Status().Update(ctx, cr), errUpdateStatus)
		}

		// delete the network Instance from the tree
		delete(r.iptree, crName)

		// We've successfully delete our resource (if necessary) and
		// removed our finalizer. If we assume we were the only controller that
		// added a finalizer to this resource then it should no longer exist and
		// thus there is no point trying to update its status.
		log.Debug("Successfully deleted resource")
		return reconcile.Result{Requeue: false}, nil
	}

	if err := r.managed.AddFinalizer(ctx, cr); err != nil {
		// If this is the first time we encounter this issue we'll be requeued
		// implicitly when we update our status with the new error condition. If
		// not, we requeue explicitly, which will trigger backoff.
		record.Event(cr, event.Warning(reasonCannotAddFInalizer, err))
		log.Debug("Cannot add finalizer", "error", err)
		cr.SetConditions(nddv1.ReconcileError(err), ipamv1alpha1.NotReady())
		return reconcile.Result{Requeue: true}, errors.Wrap(r.client.Status().Update(ctx, cr), errUpdateStatus)
	}

	if err := cr.InitializeResource(); err != nil {
		record.Event(cr, event.Warning(reasonCannotInitialize, err))
		log.Debug("Cannot initialize", "error", err)
		cr.SetConditions(nddv1.ReconcileError(err), ipamv1alpha1.NotReady())
		return reconcile.Result{Requeue: true}, errors.Wrap(r.client.Status().Update(ctx, cr), errUpdateStatus)
	}

	if err := r.handleAppLogic(ctx, cr, crName); err != nil {
		record.Event(cr, event.Warning(reasonAppLogicFailed, err))
		log.Debug("handle applogic failed", "error", err)
		cr.SetConditions(nddv1.ReconcileError(err), ipamv1alpha1.NotReady())
		return reconcile.Result{RequeueAfter: shortWait}, errors.Wrap(r.client.Status().Update(ctx, cr), errUpdateStatus)
	}

	if err := r.GarbageCollection(ctx, cr, crName); err != nil {
		record.Event(cr, event.Warning(reasonCannotGarbageCollect, err))
		log.Debug("Cannot perform garbage collection", "error", err)
		cr.SetConditions(nddv1.ReconcileError(err), ipamv1alpha1.NotReady())
		return reconcile.Result{Requeue: true}, errors.Wrap(r.client.Status().Update(ctx, cr), errUpdateStatus)
	}

	cr.SetConditions(nddv1.ReconcileSuccess(), ipamv1alpha1.Ready())
	// we don't need to requeue for ipam
	return reconcile.Result{RequeueAfter: reconcileTimeout}, errors.Wrap(r.client.Status().Update(ctx, cr), errUpdateStatus)
}

func (r *Reconciler) handleAppLogic(ctx context.Context, cr ipamv1alpha1.In, crName string) error {
	ipamname := strings.Join([]string{cr.GetOrganizationName(), cr.GetIpamName()}, ".")

	// get the deployment
	ipam := r.newIpam()
	if err := r.client.Get(ctx, types.NamespacedName{
		Namespace: cr.GetNamespace(),
		Name:      ipamname}, ipam); err != nil {
		// can happen when the deployment is not found
		cr.SetStatus("down")
		cr.SetReason("organization not found")
		return errors.Wrap(err, "organization not found")
	}
	if ipam.GetCondition(ipamv1alpha1.ConditionKindReady).Status != corev1.ConditionTrue {
		cr.SetStatus("down")
		cr.SetReason("organization not ready")
		return errors.New("organization not ready")
	}
	// ipam found

	// only initialize the iptree unless there is an ipam
	if _, ok := r.iptree[crName]; !ok {
		r.iptree[crName] = table.NewRouteTable()
	}

	if err := r.handleStatus(ctx, cr, ipam); err != nil {
		return err
	}

	cr.SetOrganizationName(cr.GetOrganizationName())
	cr.SetIpamName(cr.GetIpamName())
	cr.SetNetworkInstanceName(cr.GetNetworkInstanceName())
	return nil
}

func (r *Reconciler) handleStatus(ctx context.Context, cr ipamv1alpha1.In, ipam ipamv1alpha1.Ip) error {
	r.log.Debug("handle networkinstance status", "ipam admin status", ipam.GetAdminState(), "ipam status", ipam.GetStatus())
	if ipam.GetAdminState() == "disable" {
		cr.SetStatus("down")
		cr.SetReason("parent status down")
	} else {
		if cr.GetAdminState() == "disable" {
			cr.SetStatus("down")
			cr.SetReason("admin disabled")
		} else {
			cr.SetStatus("up")
			cr.SetReason("")
		}
	}
	return nil
}

func (r *Reconciler) GarbageCollection(ctx context.Context, cr ipamv1alpha1.In, crName string) error {
	log := r.log.WithValues("function", "garbageCollection", "Name", crName)

	// get all allocations
	alloc := r.newAllocList()
	if err := r.client.List(ctx, alloc); err != nil {
		log.Debug("Cannot get allocations", "error", err)
		return err
	}

	// check if allocations dont have IpPrefix allocated
	// check if allocations match with allocated tree -> TBD
	// -> alloc found in tree -> ok
	// -> alloc not found in tree -> assign in tree
	// we keep track of the allocated Prefixeses to compare in a second stage
	allocIpPrefixes := make([]*string, 0)
	for _, alloc := range alloc.GetAllocs() {
		// only garbage collect if the networkinstance & ipam matches
		if alloc.GetNetworkInstanceName() == cr.GetName() && alloc.GetIpamName() == cr.GetIpamName() {
			allocIpPrefix, allocIpPrefixFound := alloc.HasIpPrefix()
			if !allocIpPrefixFound {
				log.Debug("Alloc", "NetworkInstance", cr.GetName(), "Name", alloc.GetName(), "IpPrefix found", allocIpPrefixFound)
			} else {
				log.Debug("Alloc", "NetworkInstance", cr.GetName(), "Name", alloc.GetName(), "IpPrefix found", allocIpPrefixFound, "allocIpPrefix", allocIpPrefix)
			}

			var prefix string

			// the selector is used in the tree to find the entry in the tree
			// we use all the keys in the source-tag and selector for the search
			fullselector := labels.NewSelector()
			l := make(map[string]string)
			for key, val := range alloc.GetSelector() {
				req, err := labels.NewRequirement(key, selection.In, []string{val})
				if err != nil {
					log.Debug("wrong object", "Error", err)
					return err
				}
				fullselector = fullselector.Add(*req)
				l[key] = val
			}
			for key, val := range alloc.GetSourceTag() {
				req, err := labels.NewRequirement(key, selection.In, []string{val})
				if err != nil {
					log.Debug("wrong object", "Error", err)
					return err
				}
				fullselector = fullselector.Add(*req)
				l[key] = val
			}

			routes := r.iptree[crName].GetByLabel(fullselector)
			if len(routes) == 0 {
				// allocate prefix
				log.Debug("Query not found, allocate a prefix")

				// via selector perform allocation
				selector := labels.NewSelector()
				for key, val := range alloc.GetSelector() {
					req, err := labels.NewRequirement(key, selection.In, []string{val})
					if err != nil {
						log.Debug("wrong object", "Error", err)
						return err
					}
					selector = selector.Add(*req)
				}

				routes := r.iptree[crName].GetByLabel(selector)

				// required during startup when not everything is initialized
				// we break and the reconciliation will take care
				if len(routes) == 0 {
					break
				}

				prefixLength, err := getPrefixLength(alloc, cr)
				if err != nil {
					return errors.Wrap(err, "prefix Length not properly configured")
				}

				// allocation strategy -> we take the first prefix for now
				a, ok := r.iptree[crName].FindFreePrefix(routes[0].IPPrefix(), uint8(prefixLength))
				if !ok {
					log.Debug("allocation failed")
					return errors.New("allocation failed")
				}

				route := table.NewRoute(a)
				route.UpdateLabel(l)
				r.iptree[crName].Add(route)

				prefix = route.String()
			} else {
				if len(routes) > 1 {
					// this should never happen since the laels should provide uniqueness
					log.Debug("strange situation, route in tree found multiple times", "ases", routes)
				}
				prefix = routes[0].IPPrefix().String()
				switch {
				case !allocIpPrefixFound:
					log.Debug("strange situation, Prefix found in tree but alloc IP Prefix not assigned")
					alloc.SetIpPrefix(prefix)
					if err := r.client.Status().Update(ctx, alloc); err != nil {
						log.Debug("updating alloc status", "error", err)
					}
				case allocIpPrefixFound && prefix != allocIpPrefix:
					log.Debug("strange situation, Prefix found in tree but alloc AS had different IP Prefix", "tree IP Prefix", prefix, "alloc IP Prefix", allocIpPrefix)
					// can happen during init, prefix to be added
					alloc.SetIpPrefix(prefix)
					if err := r.client.Status().Update(ctx, alloc); err != nil {
						log.Debug("updating alloc status", "error", err)
					}
				default:
					// do nothing, all ok
				}
			}
			allocIpPrefixes = append(allocIpPrefixes, &prefix)
		}
	}
	// TODO based on the allocated Prefixes we collected, we can validate if the
	// tree had assigned other allocations, which dont have an alloc object

	// TBD -> how to get allocated prefixes
	/*
		found := false
		for _, treePrefix := range r.iptree[nitreename].GetAllocated() {
			for _, allocIpPrefix := range allocIpPrefixes {
				if *allocIpPrefix == treePrefix {
					found = true
					break
				}
			}
			if !found {
				log.Debug("prefix found in tree, but no alloc found -> deallocate from tree", "prefix", treePrefix)
				r.iptree[nitreename].DeAllocate(treePrefix)
			}
		}
	*/
	// always update the status field in the aspool with the latest info
	// TBD cr.UpdateAs(allocAses)

	// DUMMY
	for _, allocIpPrefix := range allocIpPrefixes {
		log.Debug("Allocated prefixes", "Prefix", allocIpPrefix)
	}

	return nil
}

func getPrefixLength(cr ipamv1alpha1.Aa, ni ipamv1alpha1.In) (uint32, error) {
	//prefixLength := cr.GetPrefixLength()
	//if prefixLength == nil {
	purpose := cr.GetSelector()[ipamv1alpha1.KeyPurpose]
	af := cr.GetSelector()[ipamv1alpha1.KeyAddressFamily]
	prefixLength := ni.GetDefaultPrefixLength(purpose, af)
	if prefixLength == nil {
		return 0, errors.New("default prefix length not configured prop")
	}
	//}
	return *prefixLength, nil
}
