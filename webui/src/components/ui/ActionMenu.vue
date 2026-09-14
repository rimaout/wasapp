<script>
import IconButton from './IconButton.vue';
import ConfirmBar from './ConfirmBar.vue';

/**
 * ActionMenu — a generic floating menu with a trigger button and a list of items.
 *
 * Each item is supplied via the `items` prop, and there are 3 types:
 *
 *   1. Simple         { id, label, icon }
 *        Picking it emits `select(item)`.
 *
 *   2. Confirm        { id, label, icon, danger: true, dangerText }
 *        Picking it morphs into a built-in confirmation screen; confirming emits
 *        `select(item)`.
 *
 *   3. Action screen  { id, label, icon, component, props }
 *        Picking it morphs into `component`, rendering it with `props` bound via
 *        v-bind. The screen reports completion via `@done` / `@cancel`, which are
 *        forwarded as the `done` event.
 *
 * To add a new morphing action you only write a screen component (or reuse an
 * existing one) and add one item object in the caller — ActionMenu is untouched.
 *
 * Example items:
 *   { id:'reply',   label:'Reply',   icon:'corner-up-left' }
 *   { id:'delete',  label:'Delete',  icon:'trash-2', danger:true, dangerText:'Are you sure you want to delete this message?' }
 *   { id:'forward', label:'Forward', icon:'share', component: ForwardActionScreen, props:{ message: msg } }
 */
export default {
	components: { IconButton, ConfirmBar },
	props: {
		items:                  { type: Array,   required: true        }, // array of menu items (simple / confirm / action screen)
		placement:              { type: String,  default: 'down-right' }, // 'down-right' | 'down-left' | 'top-right' | 'top-left' used to position the menu relative to the trigger button
		triggerButtonMode:      { type: String,  default: 'visible'    }, // 'visible' | 'hover' used to control when the trigger button is shown (hover mode is used for message menus that reveal the trigger button on hover)
		isHovering:             { type: Boolean, default: false        }, // whether the parent is hovering (used with triggerButtonMode 'hover' to reveal the trigger button)
		triggerButtonIcon:      { type: String,  required: true        }, // feather icon name for the trigger button (e.g. 'more-vertical')
		triggerButtonFilled:    { type: Boolean, default: false        }, // filled trigger background
		triggerButtonSize:      { type: String,  default: 'default'    }, // 'default' | 'small'
	},

	// Events emitted to parent component
	//  'select' - Fired when a normal item is clicked or a danger item is confirmed
	//  'done':  - Fired when a sub-component screen finishes its task
	emits: ['select', 'done'],

	data() {
		return {
			open: false,
			activeAction: null,						// The currently active action screen (if any)
			confirmItem: null,						// The currently active danger item (if any)
			menuStyle: { visibility: 'hidden' },
		};
	},

	computed: {
		showTrigger() {
			return this.triggerButtonMode !== 'hover' || this.isHovering || this.open;
		},
	},

	watch: {
		open:          { handler: 'repositionMenu', flush: 'post' },
		activeAction:  { handler: 'repositionMenu', flush: 'post' },
		confirmItem:   { handler: 'repositionMenu', flush: 'post' },
	},

	methods: {
		repositionMenu() {
			if (this.open) this.positionMenu();
		},
		closeAll() {
			// Close the menu and reset any active action or confirmation state.
			this.open = false;
			this.activeAction = null;
			this.confirmItem = null;
		},
		toggle() {
			this.open = !this.open;
		},
		choose(item) {
			// Items with a `component` morph into that action screen in place.
			if (item.component) {
				this.activeAction = { item };
				this.positionMenu();
				return;
			}
			// Danger items with a `dangerText` morph into the confirmation screen.
			if (item.danger && item.dangerText) {
				this.confirmItem = item;
				this.positionMenu();
				return;
			}
			this.closeAll();
			this.$emit('select', item);
		},
		onConfirm() {
			// Emit the selected danger item after confirmation.
			const item = this.confirmItem;
			this.closeAll();
			this.$emit('select', item);
		},
		onActionDone(payload) {
			// Emit the `done` event from the active action screen, forwarding any payload.
			this.$emit('done', payload);
			this.closeAll();
		},
		positionMenu() {
			// Safety check: check if the trigger and menu elements are available before proceeding with positioning
			if (!this.$refs.trigger || !this.$refs.menu) return;

			// Get measurements
			const triggerRect = this.$refs.trigger.getBoundingClientRect();
			const menuWidth = this.$refs.menu.offsetWidth;
			const menuHeight = this.$refs.menu.offsetHeight;

			// ----------------------------------------------------
			// STEP 1: Calculate Left Position (Horizontal)
			// ----------------------------------------------------
			let left = triggerRect.left; // Default: align left edges

			if (this.placement.includes('right')) {
				left = triggerRect.right - menuWidth; // Align right edges
			}

			// Keep menu inside screen bounds (8px margin)
			const minLeft = 8;
			const maxLeft = window.innerWidth - menuWidth - 8;

			if (left < minLeft) left = minLeft; // Don't overflow left edge
			if (left > maxLeft) left = maxLeft; // Don't overflow right edge

			// ----------------------------------------------------
			// STEP 2: Decide Open Direction (Up or Down)
			// ----------------------------------------------------
			const spaceBelow = window.innerHeight - triggerRect.bottom;
			const spaceAbove = triggerRect.top;
			const preferTop = this.placement.includes('top');

			let openUp = preferTop;

			// Flip direction if space is too tight (needs menu height + 8px gap)
			if (preferTop && spaceAbove < menuHeight + 8) {
				openUp = false; // Force open downwards
			}
			if (!preferTop && spaceBelow < menuHeight + 8) {
				openUp = true;  // Force open upwards
			}

			// ----------------------------------------------------
			// STEP 3: Apply CSS Style Object
			// ----------------------------------------------------
			const style = {
				visibility: 'visible',
				left: left + 'px'
			};

			if (openUp) {
				style.bottom = (window.innerHeight - triggerRect.top + 6) + 'px';
			} else {
				style.top = (triggerRect.bottom + 6) + 'px';
			}

			this.menuStyle = style;
		}
	},
};
</script>

<template>
	<div class="action-menu">
		<!-- Trigger button -->
		<div ref="trigger" class="action-menu-trigger" :class="{ hidden: !showTrigger }">
			<IconButton :icon="triggerButtonIcon" :filled="triggerButtonFilled" :size="triggerButtonSize" @click="toggle" />
		</div>

		<!-- Floating menu (note: `Teleport to body` is used to avoid clipping by parent containers) -->
		<Teleport to="body">
			<!-- Backdrop to close the menu when clicking outside -->
			<div v-if="open" class="menu-backdrop" @click="closeAll"></div>

			<!-- Menu content -->
			<div v-if="open" ref="menu" class="menu" :style="menuStyle">
				<!-- Render the active action screen if one is selected -->
				<component
					v-if="activeAction"
					:is="activeAction.item.component"
					v-bind="activeAction.item.props"
					@done="onActionDone"
					@cancel="closeAll"
				/>

				<!-- Render the confirmation screen if a danger item is selected (and no active action screen is selected) -->
				<div v-else-if="confirmItem" class="confirm-screen">
					<div class="danger-text">{{ confirmItem.dangerText }}</div>
					<ConfirmBar danger confirm-text="Confirm" confirm-icon="check" cancel-label="Cancel" @cancel="closeAll" @confirm="onConfirm" />
				</div>

				<!-- Render the list of menu items if no action screen or confirmation screen is active -->
				<template v-else>
					<template v-for="item in items" :key="item.id">
						<div v-if="item.danger" class="menu-divider"></div>

						<!-- Render a button for each menu item, applying danger styling if applicable -->
						<button type="button" class="menu-item" :class="{ danger: item.danger }" @click="choose(item)">
							<svg class="feather menu-icon"><use :href="'/feather-sprite-v4.29.0.svg#' + item.icon"/></svg>
							<span>{{ item.label }}</span>
						</button>

					</template>
				</template>
			</div>
		</Teleport>
	</div>
</template>

<style scoped>
.action-menu {
	display: flex;
	flex-shrink: 0;
}

.action-menu-trigger {
	display: flex;
	transition: opacity 0.15s ease;
}

.action-menu-trigger.hidden {
	opacity: 0;
	pointer-events: none;
}

.menu-backdrop {
	position: fixed;
	inset: 0;
	z-index: 200;
}

.menu {
	position: fixed;
	z-index: 201;
	min-width: 200px;
	background: var(--theme-bg-light);
	border: 1px solid var(--theme-border);
	border-radius: 12px;
	box-shadow: 0 8px 24px rgba(0, 0, 0, 0.4);
	padding: 6px;
}

.menu-item {
	display: flex;
	align-items: center;
	gap: 12px;
	width: 100%;
	background: none;
	border: none;
	border-radius: 8px;
	padding: 10px 12px;
	color: var(--theme-fg);
	font-size: 0.95rem;
	cursor: pointer;
	text-align: left;
	white-space: nowrap;
}

.menu-item:hover {
	background-color: var(--theme-overlay);
}

.menu-icon {
	width: 18px;
	height: 18px;
	color: var(--theme-fg-dark);
	flex-shrink: 0;
}

.menu-divider {
	height: 1px;
	background: var(--theme-border);
	margin: 4px 8px;
}

.menu-item.danger {
	color: var(--theme-red);
}

.menu-item.danger .menu-icon {
	color: var(--theme-red);
}

.menu-item.danger:hover {
	background-color: rgba(247, 118, 142, 0.12);
}

.confirm-screen {
	display: flex;
	flex-direction: column;
	gap: 12px;
	padding: 8px;
	width: 260px;
}

.danger-text {
	color: var(--theme-fg);
	font-size: 0.95rem;
}
</style>
