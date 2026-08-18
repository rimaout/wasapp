<script setup>
import { ref, computed, watch } from 'vue';
import IconButton from './IconButton.vue';
import PopupActionScreen from './PopupActionScreen.vue';
import ConfirmBar from './ConfirmBar.vue';

// Generalized floating action menu: a trigger button + a list of items.
// Items can be simple (emit `select`), morph into an action screen (`action`),
// or morph into a confirmation screen (`danger` + `dangerText`).
const props = defineProps({
	items: { type: Array, required: true },
	placement: { type: String, default: 'down-right' }, // 'down-right' | 'down-left' | 'top-right' | 'top-left'
	trigger: { type: String, default: 'visible' }, // 'visible' | 'hover'
	visible: { type: Boolean, default: false },
	icon: { type: String, required: true },
	label: { type: String, default: '' },
	filled: { type: Boolean, default: false },
	size: { type: String, default: 'default' },
	target: { type: Object, default: null },
	initialName: { type: String, default: '' },
});
const emit = defineEmits(['select', 'done']);

const open = ref(false);
const activeAction = ref(null);
const confirmItem = ref(null);
const trigger = ref(null);
const menu = ref(null);
const menuStyle = ref({ visibility: 'hidden' });

const showTrigger = computed(() => props.trigger !== 'hover' || props.visible || open.value);

function closeAll() {
	open.value = false;
	activeAction.value = null;
	confirmItem.value = null;
}

function toggle() {
	open.value = !open.value;
}

function choose(item) {
	// Items with an `action` morph into the action screen in place.
	if (item.action) {
		activeAction.value = { mode: item.action, title: item.label };
		positionMenu();
		return;
	}
	// Danger items with a `dangerText` morph into the confirmation screen.
	if (item.danger && item.dangerText) {
		confirmItem.value = item;
		positionMenu();
		return;
	}
	closeAll();
	emit('select', item);
}

function onConfirm() {
	const item = confirmItem.value;
	closeAll();
	emit('select', item);
}

function onActionDone(payload) {
	emit('done', payload);
	closeAll();
}

function positionMenu() {
	if (!trigger.value || !menu.value) return;
	const t = trigger.value.getBoundingClientRect();
	const mw = menu.value.offsetWidth;
	const mh = menu.value.offsetHeight;

	const alignRight = props.placement.includes('right');
	const preferTop = props.placement.includes('top');

	let left = alignRight ? t.right - mw : t.left;
	left = Math.max(8, Math.min(left, window.innerWidth - mw - 8));

	const style = { left: left + 'px', visibility: 'visible' };

	// Flip vertically if the preferred direction does not fit.
	const spaceBelow = window.innerHeight - t.bottom;
	const spaceAbove = t.top;
	let openUp = preferTop;
	if (preferTop && spaceAbove < mh + 8) openUp = false;
	if (!preferTop && spaceBelow < mh + 8) openUp = true;

	if (openUp) {
		style.bottom = (window.innerHeight - t.top + 6) + 'px';
	} else {
		style.top = (t.bottom + 6) + 'px';
	}
	menuStyle.value = style;
}

// Position (or reposition on morph) after the menu has been rendered.
watch([open, activeAction, confirmItem], () => {
	if (open.value) positionMenu();
}, { flush: 'post' });
</script>

<template>
	<div class="action-menu">
		<div ref="trigger" class="action-menu-trigger" :class="{ hidden: !showTrigger }">
			<IconButton :icon="icon" :label="label" :filled="filled" :size="size" @click="toggle" />
		</div>

		<Teleport to="body">
			<div v-if="open" class="menu-backdrop" @click="closeAll"></div>
			<div v-if="open" ref="menu" class="menu" :style="menuStyle">
				<PopupActionScreen
					v-if="activeAction"
					:mode="activeAction.mode"
					:title="activeAction.title"
					:target="target"
					:initial-name="initialName"
					@done="onActionDone"
					@cancel="closeAll"
				/>

				<div v-else-if="confirmItem" class="confirm-screen">
					<div class="danger-text">{{ confirmItem.dangerText }}</div>
					<ConfirmBar danger confirm-text="Confirm" confirm-icon="check" cancel-label="Cancel" @cancel="closeAll" @confirm="onConfirm" />
				</div>

				<template v-else>
					<template v-for="item in items" :key="item.id">
						<div v-if="item.danger" class="menu-divider"></div>
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
	background: var(--tn-bg-light);
	border: 1px solid var(--tn-border);
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
	color: var(--tn-fg);
	font-size: 0.95rem;
	cursor: pointer;
	text-align: left;
	white-space: nowrap;
}

.menu-item:hover {
	background-color: rgba(255, 255, 255, 0.08);
}

.menu-icon {
	width: 18px;
	height: 18px;
	color: var(--tn-fg-dark);
	flex-shrink: 0;
}

.menu-divider {
	height: 1px;
	background: var(--tn-border);
	margin: 4px 8px;
}

.menu-item.danger {
	color: var(--tn-red);
}

.menu-item.danger .menu-icon {
	color: var(--tn-red);
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
	color: var(--tn-fg);
	font-size: 0.95rem;
}
</style>
