package commonconstants

/**
* Message Broker Events
**/
const (
	// example
	ExampleCreatedEvent = "example.created"

	// Member Events
	MemberSignedUpEvent = "member.signedup"       // when user creates account
	PasswordResetEvent  = "member.password_reset" // when password reset is requested

	// Build events
	BuildCreatedEvent   = "build.created"   // when build is first created (draft)
	BuildPublishedEvent = "build.published" // when build is made public
	BuildUpdatedEvent   = "build.updated"   // when published build is edited
	BuildDeletedEvent   = "build.deleted"   // when build is deleted
	BuildRatedEvent     = "build.rated"     // when someone rates a build)
)

/**
* Message Broker Event Payloads
**/

/**
* MemberSignedUpEventPayload
*
* Published by auth-service.
* Consumed by:
* - notification-service
* - analytics-service
**/
type MemberSignedUpEventPayload struct {
	UserID     string `json:"userId"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	SignedUpAt string `json:"signedUpAt"`
}
