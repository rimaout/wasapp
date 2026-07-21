package database

import (
	"database/sql"
	"fmt"
)

func createTriggers(db *sql.DB) error {
	// Constraint 4.1: Recipient Exclusion for Senders
	// The recipient of a delivery/read status cannot be the sender of the message.
	triggerSenderExclusion := `CREATE TRIGGER trg_recv_status_no_sender
		BEFORE INSERT ON receiver_statuses
		FOR EACH ROW
		BEGIN
			SELECT RAISE(ABORT, 'receiver cannot be the sender of the message')
			WHERE NEW.user_id = (SELECT sender_id FROM messages WHERE id = NEW.message_id);
		END;`
	_, err := db.Exec(triggerSenderExclusion)
	if err != nil {
		return fmt.Errorf("error creating receiver_statuses sender exclusion trigger: %w", err)
	}

	// Constraint 4.2: Message Timeline Coherence (send_time <= recv_time)
	// Ensures recv_time is at or after the message's send_time (cross-table check).
	triggerTimelineInsert := `CREATE TRIGGER trg_recv_status_timeline_insert
		BEFORE INSERT ON receiver_statuses
		FOR EACH ROW
		BEGIN
			SELECT RAISE(ABORT, 'recv_time must be at or after send_time')
			WHERE NEW.recv_time IS NOT NULL
			  AND NEW.recv_time < (SELECT send_time FROM messages WHERE id = NEW.message_id);
		END;`
	_, err = db.Exec(triggerTimelineInsert)
	if err != nil {
		return fmt.Errorf("error creating receiver_statuses timeline insert trigger: %w", err)
	}

	// Constraint 4.2: Message Timeline Coherence (send_time <= recv_time)
	// Same check for UPDATE — prevents setting recv_time before send_time.
	triggerTimelineUpdate := `CREATE TRIGGER trg_recv_status_timeline_update
		BEFORE UPDATE ON receiver_statuses
		FOR EACH ROW
		WHEN NEW.recv_time IS NOT NULL AND NEW.recv_time < (SELECT send_time FROM messages WHERE id = NEW.message_id)
		BEGIN
			SELECT RAISE(ABORT, 'recv_time must be at or after send_time');
		END;`
	_, err = db.Exec(triggerTimelineUpdate)
	if err != nil {
		return fmt.Errorf("error creating receiver_statuses timeline update trigger: %w", err)
	}

	// Constraint 3.3: Interaction Forbidden on Deleted Messages
	// Blocks reactions (INSERT/UPDATE/DELETE) on messages that have been logically deleted.
	triggerReactionInsert := `CREATE TRIGGER trg_react_no_deleted_insert
		BEFORE INSERT ON reactions
		FOR EACH ROW
		WHEN (SELECT is_deleted FROM messages WHERE id = NEW.message_id) = 1
		BEGIN
			SELECT RAISE(ABORT, 'cannot react to a deleted message');
		END;`
	_, err = db.Exec(triggerReactionInsert)
	if err != nil {
		return fmt.Errorf("error creating reactions deleted msg insert trigger: %w", err)
	}

	triggerReactionUpdate := `CREATE TRIGGER trg_react_no_deleted_update
		BEFORE UPDATE ON reactions
		FOR EACH ROW
		WHEN (SELECT is_deleted FROM messages WHERE id = NEW.message_id) = 1
		BEGIN
			SELECT RAISE(ABORT, 'cannot update reaction on a deleted message');
		END;`
	_, err = db.Exec(triggerReactionUpdate)
	if err != nil {
		return fmt.Errorf("error creating reactions deleted msg update trigger: %w", err)
	}

	triggerReactionDelete := `CREATE TRIGGER trg_react_no_deleted_delete
		BEFORE DELETE ON reactions
		FOR EACH ROW
		WHEN (SELECT is_deleted FROM messages WHERE id = OLD.message_id) = 1
		BEGIN
			SELECT RAISE(ABORT, 'cannot remove reaction from a deleted message');
		END;`
	_, err = db.Exec(triggerReactionDelete)
	if err != nil {
		return fmt.Errorf("error creating reactions deleted msg delete trigger: %w", err)
	}

	// Constraint 4.3: Membership Requirement for Delivery Status
	// A delivery/read status can only be recorded for users who are active members of the chat.
	triggerRecvMembership := `CREATE TRIGGER trg_recv_status_membership
		BEFORE INSERT ON receiver_statuses
		FOR EACH ROW
		BEGIN
			SELECT RAISE(ABORT, 'user is not an active member of the chat')
			WHERE NOT EXISTS (
				SELECT 1 FROM members m
				JOIN messages msg ON m.chat_id = msg.chat_id
				WHERE msg.id = NEW.message_id
				  AND m.user_id = NEW.user_id
				  AND m.group_leave_time IS NULL
			);
		END;`
	_, err = db.Exec(triggerRecvMembership)
	if err != nil {
		return fmt.Errorf("error creating receiver_statuses membership trigger: %w", err)
	}

	return nil
}
