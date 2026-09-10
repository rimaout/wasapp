<script>
/**
 * IconButton — a reusable button component that displays a Feather icon and emits a custom click event.
 *
 * Props:
 *  - icon:     The name of the Feather icon to display.
 *  - size:     The size variant of the button ('default or 'small').
 *  - disabled: Disables button clicks when true.
 *  - filled:   Gives the button a background fill when true.
 *
 * Events:
 *  - click: Emitted when the button is clicked.
 */
export default {
	props: {
		icon:     { type: String,  required: true     },
		size:     { type: String,  default: 'default' },
		disabled: { type: Boolean, default: false     },
		filled:   { type: Boolean, default: false     },
	},
	emits: ['click'],
};
</script>

<template>
	<!-- Button with dynamic CSS classes (:class), native disabled state, and a custom click event -->
	<button
		type="button"
		class="icon-btn"
		:class="[size, { filled: filled }]"
		:disabled="disabled"
		@click="$emit('click')"
	>
		<svg class="feather icon-btn-icon"><use :href="'/feather-sprite-v4.29.0.svg#' + icon"/></svg>
	</button>
	<!-- Note: `:class="[size, { filled: filled }]"` dynamically applies the size class and conditionally adds the 'filled' class based on the if `filled` prop is true -->
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
