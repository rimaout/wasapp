<script setup>

// Define the inputs (props) passed from a parent component into this button
defineProps({
	icon:     { type: String,  required: true     }, // Feather icon name
	size:     { type: String,  default: 'default' }, // Size variant: 'default' or 'small'
	disabled: { type: Boolean, default: false     }, // Disables button clicks when true
	filled:   { type: Boolean, default: false     }, // Gives button a background fill when true
});

// Define the custom 'click' event sent back up to the parent component
const emit = defineEmits(['click']);

</script>

<template>
	<!-- Button with dynamic CSS classes (:class), native disabled state, and a custom click event -->
	<button
		type="button"
		class="icon-btn"
		:class="[size, { filled: filled }]"
		:disabled="disabled"
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
	color: var(--theme-fg-dark);
	flex-shrink: 0;
}

.icon-btn:hover {
	background-color: var(--theme-overlay);
}

.icon-btn:active {
	background-color: rgba(255, 255, 255, 0.14);
}

.icon-btn.filled {
	background-color: var(--theme-overlay);
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
