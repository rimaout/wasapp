<script>
/**
 * ConfirmBar — a two-button footer: a square X (cancel) plus a confirm button
 * that fills the remaining width. Default styling is a green confirm; setting
 * `danger` swaps it to a red confirm (used inside confirm screens).
 *
 * @property {string}  [confirmText='']        - label on the confirm button.
 * @property {string}  [confirmIcon='']        - optional feather icon name on the confirm button.
 * @property {boolean} [iconRight=false]       - put the confirm icon after the label.
 * @property {boolean} [confirmDisabled=false] - disable the confirm button.
 * @property {boolean} [cancelDisabled=false]  - disable the X (cancel) button.
 * @property {string}  [cancelLabel='Cancel']  - accessible label/tooltip for the X button.
 * @property {boolean} [danger=false]          - red confirm styling (vs. default green).
 */
export default {
	props: {
		confirmText:     { type: String,  default: ''       },
		confirmIcon:     { type: String,  default: ''       },
		iconRight:       { type: Boolean, default: false    },
		confirmDisabled: { type: Boolean, default: false    },
		cancelDisabled:  { type: Boolean, default: false    },
		cancelLabel:     { type: String,  default: 'Cancel' },
		danger:          { type: Boolean, default: false    },
	},

	/**
	 * Events:
	 *   cancel  - fired when the X button is clicked.
	 *   confirm - fired when the confirm button is clicked.
	 */
	emits: ['cancel', 'confirm'],
};
</script>

<template>
	<div class="confirm-bar" :class="{ danger }">
		<button type="button" class="footer-x" :disabled="cancelDisabled" :aria-label="cancelLabel" :title="cancelLabel" @click="$emit('cancel')">
			<svg class="feather footer-x-icon"><use href="/feather-sprite-v4.29.0.svg#x"/></svg>
		</button>
		<button type="button" class="confirm-btn" :disabled="confirmDisabled" :aria-label="confirmText || undefined" :title="confirmText || undefined" @click="$emit('confirm')">
			<svg v-if="confirmIcon && !iconRight" class="feather confirm-icon"><use :href="'/feather-sprite-v4.29.0.svg#' + confirmIcon"/></svg>
			<span v-if="confirmText">{{ confirmText }}</span>
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
