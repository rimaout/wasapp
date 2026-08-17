<script setup>
defineProps({
	icon: { type: String, required: true },
	size: { type: String, default: 'default' },
	label: { type: String, default: '' },
	disabled: { type: Boolean, default: false },
	filled: { type: Boolean, default: false },
});
const emit = defineEmits(['click']);
</script>

<template>
	<button
		type="button"
		class="icon-btn"
		:class="[size, { filled: filled }]"
		:disabled="disabled"
		:aria-label="label || undefined"
		:title="label || undefined"
		@click="emit('click')"
	>
		<svg class="feather icon-btn-icon"><use :href="'/feather-sprite-v4.29.0.svg#' + icon"/></svg>
	</button>
</template>

<style scoped>
.icon-btn {
	display: flex;
	align-items: center;
	justify-content: center;
	background: none;
	border: none;
	border-radius: 50%;
	padding: 0;
	cursor: pointer;
	color: var(--tn-fg-dark);
	flex-shrink: 0;
}

.icon-btn:hover {
	background-color: rgba(255, 255, 255, 0.08);
}

.icon-btn:active {
	background-color: rgba(255, 255, 255, 0.14);
}

.icon-btn.filled {
	background-color: rgba(255, 255, 255, 0.08);
}

.icon-btn.filled:hover {
	background-color: rgba(255, 255, 255, 0.14);
}

.icon-btn.filled:active {
	background-color: rgba(255, 255, 255, 0.2);
}

.icon-btn:disabled {
	cursor: default;
	opacity: 0.5;
}

.icon-btn.default {
	width: 40px;
	height: 40px;
}

.icon-btn.default .icon-btn-icon {
	width: 23px;
	height: 23px;
}

.icon-btn.small {
	width: 32px;
	height: 32px;
}

.icon-btn.small .icon-btn-icon {
	width: 16px;
	height: 16px;
}
</style>
