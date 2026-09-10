<script>
/**
 * ConfirmBar — a two-button footer: a square X (cancel) plus a confirm button
 * that fills the remaining width.
 *
 * Props:
 *  - confirmText:     Label on the confirm button
 *  - confirmIcon:     Optional feather icon name on the confirm button
 *  - iconRight:       Put the confirm icon after the label (otherwise it is before the label)
 *  - confirmDisabled: Disable the confirm button (for example when a form is invalid)
 *  - cancelDisabled:  Disable the X (cancel) button
 *  - danger:          Red confirm styling (vs. default green)
 *
 * Events:
 *  - cancel:  Fired when the X button is clicked
 *  - confirm: Fired when the confirm button is clicked
 */
export default {
	props: {
		confirmText:     { type: String,  default: ''       },
		confirmIcon:     { type: String,  default: ''       },
		iconRight:       { type: Boolean, default: false    },
		confirmDisabled: { type: Boolean, default: false    },
		cancelDisabled:  { type: Boolean, default: false    },
		danger:          { type: Boolean, default: false    },
	},
	emits: ['cancel', 'confirm'],
};
</script>

<template>
	<div class="confirm-bar" :class="{ danger }">

		<!-- Cancel button (X) -->
		<button type="button" class="footer-x" :disabled="cancelDisabled" @click="$emit('cancel')">
			<svg class="feather footer-x-icon"><use href="/feather-sprite-v4.29.0.svg#x"/></svg>
		</button>

		<!-- Confirm button -->
		<button type="button" class="confirm-btn" :disabled="confirmDisabled" :title="confirmText || undefined" @click="$emit('confirm')">
			<!-- Confirm icon (optional) -->
			<svg v-if="confirmIcon && !iconRight" class="feather confirm-icon"><use :href="'/feather-sprite-v4.29.0.svg#' + confirmIcon"/></svg>

			<!-- Confirm label (optional) -->
			<span v-if="confirmText">{{ confirmText }}</span>

			<!-- Confirm icon (optional, right-aligned) -->
			<svg v-if="confirmIcon && iconRight" class="feather confirm-icon"><use :href="'/feather-sprite-v4.29.0.svg#' + confirmIcon"/></svg>
		</button>
	</div>
</template>

<style scoped>
.confirm-bar {
	display: flex;
	gap: 8px;
}

.footer-x {
	display: flex;
	align-items: center;
	justify-content: center;
	width: 40px;
	height: 40px;
	flex-shrink: 0;
	border: none;
	border-radius: 8px;
	cursor: pointer;
	background: rgba(247, 118, 142, 0.12);
	color: var(--theme-red);
}

.footer-x:hover:not(:disabled) {
	background: rgba(247, 118, 142, 0.22);
}

.footer-x:disabled {
	opacity: 0.5;
	cursor: default;
}

.footer-x-icon {
	width: 22px;
	height: 22px;
}

.confirm-btn {
	display: flex;
	align-items: center;
	justify-content: center;
	gap: 6px;
	flex: 1 1 auto;
	height: 40px;
	border: none;
	border-radius: 8px;
	cursor: pointer;
	background: var(--theme-safe-green-bg);
	color: var(--theme-green);
	font-weight: 600;
	font-size: 0.9rem;
}

.confirm-btn:hover:not(:disabled) {
	background: rgba(158, 206, 106, 0.25);
}

.confirm-btn:disabled {
	background: var(--theme-overlay);
	color: var(--theme-fg-dark);
	cursor: default;
}

.confirm-bar.danger .footer-x {
	background: var(--theme-overlay);
	color: var(--theme-fg-dark);
}

.confirm-bar.danger .footer-x:hover:not(:disabled) {
	background: rgba(255, 255, 255, 0.16);
}

.confirm-bar.danger .confirm-btn {
	background: var(--theme-danger-red-bg);
	color: var(--theme-red);
}

.confirm-bar.danger .confirm-btn:hover:not(:disabled) {
	background: rgba(247, 118, 142, 0.25);
}

.confirm-icon {
	width: 18px;
	height: 18px;
}
</style>
