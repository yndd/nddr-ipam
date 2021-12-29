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

package ipamnetworkinstanceipprefix

import (
	"context"
	"fmt"
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
	"inet.af/netaddr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const (
	finalizerName = "finalizer.ipamnetworkinstanceipprefix.ipam.nddr.yndd.io"
	//
	reconcileTimeout = 1 * time.Minute
	longWait         = 1 * time.Minute
	mediumWait       = 30 * time.Second
	shortWait        = 15 * time.Second
	veryShortWait    = 5 * time.Second

	// Errors
	errGetK8sResource = "cannot get ipamnetworkinstanceipprefix resource"
	errUpdateStatus   = "cannot update status of ipamnetworkinstanceipprefix resource"

	// events
	reasonReconcileSuccess             event.Reason = "ReconcileSuccess"
	reasonCannotDelete                 event.Reason = "CannotDeleteResource"
	reasonCannotAddFInalizer           event.Reason = "CannotAddFinalizer"
	reasonCannotDeleteFInalizer        event.Reason = "CannotDeleteFinalizer"
	reasonCannotInitialize             event.Reason = "CannotInitializeResource"
	reasonCannotGetAllocations         event.Reason = "CannotGetAllocations"
	reasonAppLogicFailed               event.Reason = "ApplogicFailed"
	reasonCannotParseIpPrefix          event.Reason = "CannotParseIpPrefix"
	reasonCannotDeleteDueToAllocations event.Reason = "CannotDeleteIpPrefixDueToExistingAllocations"
)

// ReconcilerOption is used to configure the Reconciler.
type ReconcilerOption func(*Reconciler)

// Reconciler reconciles packages.
type Reconciler struct {
	client  resource.ClientApplicator
	log     logging.Logger
	record  event.Recorder
	managed mrManaged

	newIpamNetworkInstanceIpPrefix func() ipamv1alpha1.Ipp
	newIpamNetworkInstance         func() ipamv1alpha1.In

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

func WithNewReourceFn(f func() ipamv1alpha1.Ipp) ReconcilerOption {
	return func(r *Reconciler) {
		r.newIpamNetworkInstanceIpPrefix = f
	}
}

func WithIpTree(iptree map[string]*table.RouteTable) ReconcilerOption {
	return func(r *Reconciler) {
		r.iptree = iptree
	}
}

func WithNewIpamNetworkInstanceFn(f func() ipamv1alpha1.In) ReconcilerOption {
	return func(r *Reconciler) {
		r.newIpamNetworkInstance = f
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
	name := "nddr/" + strings.ToLower(ipamv1alpha1.IpamNetworkInstanceIpPrefixGroupKind)
	fn := func() ipamv1alpha1.Ipp { return &ipamv1alpha1.IpamNetworkInstanceIpPrefix{} }
	nifn := func() ipamv1alpha1.In { return &ipamv1alpha1.IpamNetworkInstance{} }

	r := NewReconciler(mgr,
		WithLogger(nddcopts.Logger.WithValues("controller", name)),
		WithNewReourceFn(fn),
		WithIpTree(nddcopts.Iptree),
		WithNewIpamNetworkInstanceFn(nifn),
		WithRecorder(event.NewAPIRecorder(mgr.GetEventRecorderFor(name))),
	)

	return ctrl.NewControllerManagedBy(mgr).
		Named(name).
		WithOptions(o).
		For(&ipamv1alpha1.IpamNetworkInstanceIpPrefix{}).
		//Watches(&source.Kind{Type: &ipamv1alpha1.Ipam{}}, ipamHandler).
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
	log.Debug("Reconciling IpamNetworkInstanceIpPrefix", "NameSpace", req.NamespacedName)

	ctx, cancel := context.WithTimeout(ctx, reconcileTimeout)
	defer cancel()

	cr := r.newIpamNetworkInstanceIpPrefix()
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

		// check if allocations exists
		if _, ok := r.iptree[crName]; ok {
			p, err := netaddr.ParseIPPrefix(cr.GetPrefix())
			if err != nil {
				record.Event(cr, event.Warning(reasonCannotParseIpPrefix, err))
				log.Debug("Cannot parse ip prefix", "error", err)
				cr.SetConditions(nddv1.ReconcileError(err), ipamv1alpha1.NotReady())
				return reconcile.Result{Requeue: true}, errors.Wrap(r.client.Status().Update(ctx, cr), errUpdateStatus)
			}
			routes := r.iptree[crName].Children(p)

			if len(routes) > 0 {
				// We cannot delete the prefix yet due to existing allocations
				record.Event(cr, event.Warning(reasonCannotDeleteDueToAllocations, err))
				log.Debug("Cannot delete prefix due to existing allocations", "error", err)
				return reconcile.Result{RequeueAfter: shortWait}, errors.Wrap(r.client.Status().Update(ctx, cr), errUpdateStatus)
			}

			route := table.NewRoute(p)
			route.UpdateLabel(cr.GetTags())

			if _, _, err := r.iptree[crName].Delete(route); err != nil {
				log.Debug("IPPrefix deleteion failed", "prefix", p)
				cr.SetStatus("down")
				cr.SetReason(fmt.Sprintf("IPPrefix deletion failed %v", err))
				cr.SetConditions(nddv1.ReconcileError(err), ipamv1alpha1.NotReady())
				return reconcile.Result{Requeue: true}, errors.Wrap(r.client.Status().Update(ctx, cr), errUpdateStatus)
			}
		}

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

	cr.SetConditions(nddv1.ReconcileSuccess(), ipamv1alpha1.Ready())
	return reconcile.Result{}, errors.Wrap(r.client.Status().Update(ctx, cr), errUpdateStatus)
}

func (r *Reconciler) handleAppLogic(ctx context.Context, cr ipamv1alpha1.Ipp, crName string) error {
	log := r.log.WithValues("name", cr.GetName())
	log.Debug("handle application logic")
	niName := strings.Join([]string{cr.GetOrganizationName(), cr.GetIpamName(), cr.GetNetworkInstanceName()}, ".")

	// get the deployment
	ni := r.newIpamNetworkInstance()
	if err := r.client.Get(ctx, types.NamespacedName{
		Namespace: cr.GetNamespace(),
		Name:      niName}, ni); err != nil {
		// can happen when the deployment is not found
		cr.SetStatus("down")
		cr.SetReason("network instance not found")
		return errors.Wrap(err, "network instance not found")
	}
	if ni.GetCondition(ipamv1alpha1.ConditionKindReady).Status != corev1.ConditionTrue {
		cr.SetStatus("down")
		cr.SetReason("network instance not ready")
		return errors.New("network instance not ready")
	}

	if _, ok := r.iptree[crName]; !ok {
		log.Debug("Parent Routing table not ready")
		cr.SetStatus("down")
		cr.SetReason("Parent Routing table not ready")
		return errors.New("ipam not ready")
	}

	p, err := netaddr.ParseIPPrefix(cr.GetPrefix())
	if err != nil {
		log.Debug("UpdateConfig ParseIPPrefix", "Error", err)
		cr.SetStatus("down")
		cr.SetReason(fmt.Sprintf("parse ip prefix failed %v", err))
		return errors.Wrap(err, "ParseIPPrefix failed")
	}
	// we derive the address family from the prefix, to avoid exposing it to the user
	var af string
	if p.IP().Is4() {
		af = string(ipamv1alpha1.AddressFamilyIpv4)
	}
	if p.IP().Is6() {
		af = string(ipamv1alpha1.AddressFamilyIpv6)
	}
	// we add the address family in the tag/label to allow to selec the prefix on this basis
	tags := cr.GetTags()
	tags[ipamv1alpha1.KeyAddressFamily] = af
	route := table.NewRoute(p)
	route.UpdateLabel(tags)

	if err := r.iptree[crName].Add(route); err != nil {
		log.Debug("IPPrefix insertion failed", "prefix", p)
		cr.SetStatus("down")
		cr.SetReason(fmt.Sprintf("IPPrefix insertion failed %v", err))
		return errors.Wrap(err, "IPPrefix insertion failed")
	}
	log.Debug("IPPrefix insert success", "prefix", p)
	cr.SetStatus("up")
	cr.SetReason("")
	cr.SetAddressFamily(af)
	cr.SetOrganizationName(cr.GetOrganizationName())
	cr.SetIpamName(cr.GetIpamName())
	cr.SetNetworkInstanceName(cr.GetNetworkInstanceName())

	return nil

}
