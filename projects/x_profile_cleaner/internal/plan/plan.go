package plan

import (
	"fmt"
	"unicode/utf8"

	"git.alwaldend.com/alwaldend/src/projects/x_profile_cleaner/api/workerpb"
	"git.alwaldend.com/alwaldend/src/projects/x_profile_cleaner/internal/model"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

const (
	// MaxBytes reserves transport space for trusted job metadata in a 16 MiB response.
	MaxBytes = 16*1024*1024 - 4096
	// MaxOperations bounds operation allocation independently of encoded size.
	MaxOperations = 50000
)

// OperationName returns the schema-owned name of the selected payload for diagnostics.
func OperationName(value *workerpb.Operation) string {
	if value == nil {
		return ""
	}
	message := value.ProtoReflect()
	field := message.WhichOneof(message.Descriptor().Oneofs().Get(0))
	if field == nil {
		return ""
	}
	return string(field.Name())
}

// TargetID returns a referenced post/account ID; creation and profile updates have none.
func TargetID(value *workerpb.Operation) string {
	switch action := value.GetAction().(type) {
	case *workerpb.Operation_DeletePost:
		return action.DeletePost.GetPostId()
	case *workerpb.Operation_Unrepost:
		return action.Unrepost.GetPostId()
	case *workerpb.Operation_Unlike:
		return action.Unlike.GetPostId()
	case *workerpb.Operation_Unfollow:
		return action.Unfollow.GetUserId()
	case *workerpb.Operation_ReadPost:
		return action.ReadPost.GetPostId()
	default:
		return ""
	}
}

// Validate rejects malformed, unknown or excessive operation plans before publication.
func Validate(value *workerpb.Plan) error {
	if value == nil || len(value.Operations) > MaxOperations || proto.Size(value) > MaxBytes {
		return fmt.Errorf("missing or excessive operation plan")
	}
	if err := validateMessage(value.ProtoReflect()); err != nil {
		return fmt.Errorf("validate plan fields: %w", err)
	}
	for _, operation := range value.Operations {
		if operation == nil {
			return fmt.Errorf("missing operation")
		}
		id := ""
		switch action := operation.Action.(type) {
		case *workerpb.Operation_DeletePost:
			if action.DeletePost == nil {
				return fmt.Errorf("missing delete-post payload")
			}
			id = action.DeletePost.PostId
		case *workerpb.Operation_Unrepost:
			if action.Unrepost == nil {
				return fmt.Errorf("missing unrepost payload")
			}
			id = action.Unrepost.PostId
		case *workerpb.Operation_Unlike:
			if action.Unlike == nil {
				return fmt.Errorf("missing unlike payload")
			}
			id = action.Unlike.PostId
		case *workerpb.Operation_Unfollow:
			if action.Unfollow == nil {
				return fmt.Errorf("missing unfollow payload")
			}
			id = action.Unfollow.UserId
		case *workerpb.Operation_ReadPost:
			if action.ReadPost == nil {
				return fmt.Errorf("missing read-post payload")
			}
			id = action.ReadPost.PostId
		case *workerpb.Operation_CreatePost:
			if action.CreatePost == nil || action.CreatePost.Text == "" {
				return fmt.Errorf("missing create-post text")
			}
			continue
		case *workerpb.Operation_UpdateProfile:
			if action.UpdateProfile == nil || action.UpdateProfile.Description == nil {
				return fmt.Errorf("missing profile description update")
			}
			continue
		default:
			return fmt.Errorf("missing or unknown operation payload")
		}
		if !model.ValidID(id) {
			return fmt.Errorf("invalid operation target identifier")
		}
	}
	if report := value.CleanupReport; report != nil {
		if len(report.Decisions) > MaxOperations {
			return fmt.Errorf("excessive cleanup report")
		}
		for _, decision := range report.Decisions {
			if decision == nil || decision.CreatedAt == nil {
				return fmt.Errorf("missing cleanup decision or time")
			}
			if err := decision.CreatedAt.CheckValid(); err != nil {
				return fmt.Errorf("validate cleanup timestamp: %w", err)
			}
			for _, reason := range decision.Reasons {
				if reason == nil || reason.RuleId <= 0 {
					return fmt.Errorf("invalid cleanup reason")
				}
			}
		}
	}
	return nil
}

// validateMessage rejects unknown fields and invalid strings throughout the finite schema.
func validateMessage(message protoreflect.Message) error {
	if !message.IsValid() || len(message.GetUnknown()) > 0 {
		return fmt.Errorf("missing message or unknown protobuf fields")
	}
	var result error
	message.Range(func(field protoreflect.FieldDescriptor, value protoreflect.Value) bool {
		validate := func(value protoreflect.Value) bool {
			if field.Kind() == protoreflect.MessageKind {
				if err := validateMessage(value.Message()); err != nil {
					result = fmt.Errorf("validate nested message: %w", err)
					return false
				}
			}
			if field.Kind() == protoreflect.StringKind && !utf8.ValidString(value.String()) {
				result = fmt.Errorf("invalid UTF-8 plan string")
				return false
			}
			return true
		}
		if field.IsList() {
			list := value.List()
			if list.Len() > MaxOperations {
				result = fmt.Errorf("excessive repeated plan field")
				return false
			}
			for index := 0; index < list.Len(); index++ {
				if !validate(list.Get(index)) {
					return false
				}
			}
			return true
		}
		return validate(value)
	})
	if result != nil {
		return fmt.Errorf("validate protobuf message: %w", result)
	}
	return nil
}
