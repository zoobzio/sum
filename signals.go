package sum

import "github.com/zoobzio/capitan"

// Resource operation signals.
var (
	// ResourceCreated is emitted when a resource is successfully created.
	ResourceCreated = capitan.NewSignal(
		"sum.resource.created",
		"Resource was successfully created",
	)

	// ResourceListed is emitted when resources are listed.
	ResourceListed = capitan.NewSignal(
		"sum.resource.listed",
		"Resources were listed",
	)

	// ResourceRetrieved is emitted when a resource is retrieved by ID.
	ResourceRetrieved = capitan.NewSignal(
		"sum.resource.retrieved",
		"Resource was retrieved",
	)

	// ResourceUpdated is emitted when a resource is successfully updated.
	ResourceUpdated = capitan.NewSignal(
		"sum.resource.updated",
		"Resource was successfully updated",
	)

	// ResourceDeleted is emitted when a resource is successfully deleted.
	ResourceDeleted = capitan.NewSignal(
		"sum.resource.deleted",
		"Resource was successfully deleted",
	)

	// ResourceNotFound is emitted when a requested resource does not exist.
	ResourceNotFound = capitan.NewSignal(
		"sum.resource.not_found",
		"Requested resource was not found",
	)

	// ResourceError is emitted when a resource operation fails.
	ResourceError = capitan.NewSignal(
		"sum.resource.error",
		"Resource operation failed",
	)
)

// Hook lifecycle signals.
var (
	// HookExecuting is emitted when a hook pipeline begins execution.
	HookExecuting = capitan.NewSignal(
		"sum.hook.executing",
		"Hook pipeline is being executed",
	)

	// HookCompleted is emitted when a hook pipeline completes successfully.
	HookCompleted = capitan.NewSignal(
		"sum.hook.completed",
		"Hook pipeline completed successfully",
	)

	// HookFailed is emitted when a hook pipeline fails with an error.
	HookFailed = capitan.NewSignal(
		"sum.hook.failed",
		"Hook pipeline failed with error",
	)

	// HookAfterError is emitted when an after-hook fails but the request continues.
	HookAfterError = capitan.NewSignal(
		"sum.hook.after.error",
		"After hook failed but request continues",
	)

	// HookRegistered is emitted when a hook processor is registered.
	HookRegistered = capitan.NewSignal(
		"sum.hook.registered",
		"Hook processor registered for lifecycle point",
	)
)

// Field keys for resource and hook events.
var (
	// KeyResource identifies the resource type name.
	KeyResource = capitan.NewStringKey("resource")

	// KeyEntityID identifies the entity being processed.
	KeyEntityID = capitan.NewStringKey("entity_id")

	// KeyOperation identifies the operation type (create, update, delete, etc.).
	KeyOperation = capitan.NewStringKey("operation")

	// KeyCount records the number of items (e.g., in list operations).
	KeyCount = capitan.NewInt64Key("count")

	// KeyDurationMs records the operation duration in milliseconds.
	KeyDurationMs = capitan.NewInt64Key("duration_ms")

	// KeyError contains the error message when an operation fails.
	KeyError = capitan.NewStringKey("error")

	// KeyHook identifies the hook point (e.g., "beforeCreate").
	KeyHook = capitan.NewStringKey("hook")
)
