<script setup>
import { ref } from 'vue';
import IconButton from './IconButton.vue';
import PopupActionScreen from './PopupActionScreen.vue';

defineProps({
	icon: { type: String, required: true },
	label: { type: String, default: '' },
	items: { type: Array, required: true },
	direction: { type: String, default: 'down' },
	target: { type: Object, default: null },
	initialName: { type: String, default: '' },
});
const emit = defineEmits(['select', 'done']);

const open = ref(false);
const activeAction = ref(null);

function closeAll() {
	open.value = false;
	activeAction.value = null;
}

function choose(item) {
	// Items with an `action` morph the menu into the action screen in place.
	if (item.action) {
		activeAction.value = { mode: item.action, title: item.label };
		return;
	}
	closeAll();
	emit('select', item);
}

function onDone(payload) {
	emit('done', payload);
	closeAll();
}
</script>

<template>
	<div class="popup-menu">
		<IconButton :icon="icon" :label="label" @click="open = !open" />
		<template v-if="open">
			<div class="popup-backdrop" @click="closeAll"></div>
			<div class="popup-panel" :class="direction">
				<PopupActionScreen
					v-if="activeAction"
					:mode="activeAction.mode"
					:title="activeAction.title"
					:target="target"
					:initial-name="initialName"
					@done="onDone"
					@cancel="closeAll"
				/>
				<template v-else>
					<button
						v-for="item in items"
						:key="item.id"
						type="button"
						class="popup-item"
						:class="{ danger: item.danger }"
						@click="choose(item)"
					>
						<svg class="feather popup-icon"><use :href="'/feather-sprite-v4.29.0.svg#' + item.icon"/></svg>
						<span>{{ item.label }}</span>
					</button>
				</template>
			</div>
		</template>
	</div>
</template>

<style scoped>
.popup-menu {
	position: relative;
	display: flex;
}

.popup-backdrop {
	position: fixed;
	inset: 0;
	z-index: 40;
}

.popup-panel {
	position: absolute;
	top: calc(100% + 6px);
	right: 0;
	min-width: 220px;
	background: var(--tn-bg-light);
	border: 1px solid var(--tn-border);
	border-radius: 12px;
	box-shadow: 0 8px 24px rgba(0, 0, 0, 0.4);
	padding: 6px;
	z-index: 50;
}

.popup-panel.up {
	top: auto;
	bottom: calc(100% + 6px);
}

.popup-item {
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

.popup-item:hover {
	background-color: rgba(255, 255, 255, 0.08);
}

.popup-item.danger {
	color: var(--tn-red);
}

.popup-item.danger:hover {
	background-color: rgba(247, 118, 142, 0.12);
}

.popup-item .popup-icon {
	width: 18px;
	height: 18px;
	color: var(--tn-fg-dark);
	flex-shrink: 0;
}

.popup-item.danger .popup-icon {
	color: var(--tn-red);
}
</style>
