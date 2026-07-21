# Integrity Constraint Specifications

This document defines the set of integrity constraints applied to the database.

## 1) Message Constraints

>[!NOTE] ✅ 1.1 Minimum Content Requirement
>
>A message must not be empty. It must contain either text, an image, or both. Unless it is marked as deleted, is an initialization message, or is a forwarded message.
>
>$$\neg \text{isInitMessage} \land \neg \text{isDeleted} \land \neg \text{isForwardMessage} \implies \text{text} \neq \text{NULL} \lor \text{image\_id} \neq \text{NULL}$$
>
>>**Implementation:** CHECK constraint `chk_msg_minimum_content` on `messages` table in `service/database/database.go`.

>[!NOTE] ✅ 1.2 Initial Message (`isInitMessage`)
>
>A message marked as a Init Message must not contain images or text. Also it can't be logically deleted or be a forward message.
>
>$$\text{isInitMessage} \implies \text{text} = \text{NULL} \land \text{image\_id} = \text{NULL} \land \text{isDeleted} = \text{False} \land \text{isForwardMessage} = \text{False}$$
>
>>**Implementation:** CHECK constraint `chk_init_msg_fields` on `messages` table in `service/database/database.go`.

>[!NOTE] ✅ 1.3 Forwarded Message Parameters
>
>A forwarded message is represented by three fields (`isForwardMessage`, `forwarded_from_chat_id`, and `forwarded_from_msg_id`), the last two are used to trace the message back to its origin.
>
>This introduces two rules:
>
>$$\text{isForwardMessage} \iff \text{forwarded\_from\_chat\_id} \neq \text{NULL} \lor \text{forwarded\_from\_msg\_id} \neq \text{NULL}$$
>
>$$\text{forwarded\_from\_chat\_id} \neq \text{NULL} \iff \text{forwarded\_from\_msg\_id} \neq \text{NULL}$$
>
>>**Implementation:** CHECK constraint `chk_forward_pointers` on `messages` table in `service/database/database.go`.

>[!NOTE] ✅ 1.4 Content of a Forwarded Message
>
>To prevent data duplication a forwarded message record must not contain raw images or raw text directly. Instead, these contents are dynamically resolved from the original pointers.
>
>$$\text{isForwardMessage} \implies \text{text} = \text{NULL} \land \text{image\_id} = \text{NULL} \land \text{isInitMessage} = \text{False}$$
>
>>**Implementation:** CHECK constraint `chk_forward_content_empty` on `messages` table in `service/database/database.go`.

>[!NOTE] ✅ 1.5 Logical Deletion of Data (`isDeleted`)
>
>When a message is logically deleted, all direct references to user-generated data (text and image references) must be removed.
>
>$$\text{isDeleted} \implies \text{text} = \text{NULL} \land \text{image\_id} = \text{NULL}$$
>
>If the logically deleted message was a forward, the links pointing to the origin chat and original message must be removed:
>
>$$\text{isDeleted} \implies \text{isForwardMessage} = \text{False} \land \text{forwarded\_from\_chat\_id} = \text{NULL} \land \text{forwarded\_from\_msg\_id} = \text{NULL}$$
>
>Additionally, a system-generated initialization message cannot be deleted:
>
>$$\text{isDeleted} \implies \neg \text{isInitMessage}$$
>
>>**Implementation:** CHECK constraint `chk_logical_deletion_state` on `messages` table in `service/database/database.go`.

>[!NOTE] ✅ 1.6 Reply Context Coherence (`replyTo`)
>
>A message can only be submitted as a reply to another message if the target message belongs to the exact same chat.
>
>$$\text{replyTo} \neq \text{NULL} \implies \text{chat\_id} = \text{replyTo.chat\_id}$$
>
>>**Implementation:** Composite foreign key `(chat_id, reply_to_msg_id) REFERENCES messages(chat_id, id)` with `UNIQUE(chat_id, id)` on `messages` table in `service/database/database.go`.

## 2) Chat and Participant Constraints

>[!NOTE] 2.1 Uniqueness of Private Chats
>
>There can exist at most one private chat between any pair of users.
>
>$$\forall U_1, U_2, C_1 : U_{1}, U_{2} \in \text{Users} \land C_{1} \in \text{PrivateChat} \land U_{1}, U_{2} \in C_{1} \implies \neg \exists C_{2} \in \text{PrivateChat} : U_{1}, U_{2} \in C_{2}$$

>[!NOTE] 2.2 Member Limit in Private Chats
>
>A private direct chat must contain exactly two participants and individual members cannot "leave" a private chat.
>
>$$\vert \text{PrivateChat.members} \vert = 2$$

>[!NOTE] 2.3 Maximum Capacity of Group Chats
>
>A group chat cannot exceed a maximum threshold of 1024 concurrent active members. A active member is a user that did not leave the group (`leave_time = null`)

>[!NOTE] 2.4 Send a Message Only if active Member
>
>A user is authorized to post messages within a chat only if they were an active member at the time of sending.
>
>$$\text{sendTime} \ge \text{joinTime} \land (\text{leaveTime} = \text{NULL} \lor \text{sendTime} \le \text{leaveTime})$$

## Reaction Constraints

>[!NOTE] ✅ 3.1 Reaction Uniqueness per User
>
>A participant can submit at most one reaction (emoji) on any single message. Registering a new reaction on the same message automatically overwrites the previous selection.
>
>>**Implementation:** Composite primary key `(message_id, user_id)` on `reactions` table in `service/database/database.go`.

>[!NOTE] 3.2 Chat Membership Requirement for Reactions
>
>A user is only permitted to react to a message if they are an active participant of the chat containing that message.
>
>$$\text{Reaction.user\_id} \in \{ \text{members of } \text{Message.chat\_id} \}$$

>[!NOTE] ✅ 3.3 Interaction Forbidden on Deleted Messages
>
>It is strictly prohibited to insert, update, or remove reactions on messages that have been logically deleted.
>
>$$\text{Message.isDeleted} \implies \text{Reaction operations are disabled}$$
>
>>**Implementation:** Triggers `trg_react_no_deleted_insert`, `trg_react_no_deleted_update`, `trg_react_no_deleted_delete` on `reactions` table in `service/database/db-triggers.go`.

>[!NOTE] 3.4 Send a Reaction Only if active Member
>
>A user is authorized to react to messages within a chat only if they were an active member at the time of reacting.
>
>$$\text{reactTime} \ge \text{joinTime} \land (\text{leaveTime} = \text{NULL} \lor \text{reactTime} \le \text{leaveTime})$$

## 4) Message Delivery Status Constraints (`receiver_status`)

>[!NOTE] ✅ 4.1 Recipient Exclusion for Senders
>
>The recipient of a delivery or read status update cannot be the original sender of that message.
>
>$$\text{receiver\_status.user\_id} \neq \text{Message.sender\_id}$$
>
>>**Implementation:** Trigger `trg_recv_status_no_sender` on `receiver_statuses` table in `service/database/db-triggers.go`.

>[!NOTE] ✅ 4.2 Message Timeline Coherence
>
>The status transition of a message must follow chronological order: a message must be sent before it can be delivered, and it must be delivered before it can be read.
>
>$$\text{sendTime} \le \text{recvTime} \le \text{readTime}$$
>
>>**Implementation:** CHECK constraint `chk_delivery_timeline` (`recv_time <= read_time`) on `receiver_statuses` table in `service/database/database.go` + triggers `trg_recv_status_timeline_insert`, `trg_recv_status_timeline_update` (`send_time <= recv_time`) in `service/database/db-triggers.go`.

>[!NOTE] ✅ 4.3 Membership Requirement for Delivery Status
>
>A message delivery or read status update record can only be associated with users who are active members of the chat where the message belongs.
>
>$$\text{receiver\_status.user\_id} \in \{ \text{members of } \text{Message.chat\_id} \}$$
>
>>**Implementation:** Trigger `trg_recv_status_membership` on `receiver_statuses` table in `service/database/db-triggers.go`.

## 5) Security Constraints on Forwarding

>[!NOTE] 5.1 Forward Authorization (Privacy)
>
>An active user $X$ is authorized to forward a message $M$ if and only if $X$ is a registered, active member of the origin chat where message $M$ resides.
>
>$$X \in \{ \text{members of } \text{Message}(M)\text{.chat\_id} \}$$
