<script>
export default {
	props: {
		modelValue: { type: String, default: '' },
		sending: { type: Boolean, default: false },
	},
	emits: ['update:modelValue', 'send'],
	computed: {
		text: {
			get() { return this.modelValue; },
			set(v) { this.$emit('update:modelValue', v); },
		},
	},
	methods: {
		handleSend() {
			if (!this.text.trim() || this.sending) return;
			this.$emit('send');
		},
	},
};
</script>

<template>
	<div class="chat-input">
		<div class="message-input-group">
			<svg class="feather image-icon" style="width: 18px; height: 18px;">
				<use href="/feather-sprite-v4.29.0.svg#image"/>
			</svg>
			<input
				type="text"
				class="message-field"
				placeholder="Type a message..."
				v-model="text"
				@keyup.enter="handleSend"
				:disabled="sending"
			/>
			<button v-if="text.trim()" class="send-btn" @click="handleSend" :disabled="sending">
				<svg class="feather" style="width: 20px; height: 20px;"><use href="/feather-sprite-v4.29.0.svg#send"/></svg>
			</button>
		</div>
	</div>
</template>

<style scoped>
.chat-input {
	flex-shrink: 0;
	border-top: 1px solid var(--tn-border);
	padding: 10px 10px 14px;
}

.message-input-group {
	display: flex;
	align-items: center;
	background: var(--tn-bg-highlight);
	border: 1px solid var(--tn-border);
	border-radius: 999px;
	height: 52px;
	padding: 6px 6px 6px 18px;
}

.image-icon {
	color: var(--tn-comment);
	flex-shrink: 0;
	margin-right: 12px;
}

.message-field {
	flex: 1;
	background: none;
	border: none;
	outline: none;
	color: var(--tn-fg);
	font-size: 0.95rem;
	min-width: 0;
}

.message-field::placeholder {
	color: var(--tn-comment);
}

.send-btn {
	width: 38px;
	height: 38px;
	display: flex;
	align-items: center;
	justify-content: center;
	background-color: var(--tn-blue);
	border: none;
	border-radius: 50%;
	color: var(--tn-bg-darker);
	cursor: pointer;
	flex-shrink: 0;
	box-shadow: 0 0 10px rgba(122, 162, 247, 0.5);
}

.send-btn:hover {
	background-color: var(--tn-cyan);
	box-shadow: 0 0 14px rgba(125, 207, 255, 0.7);
}

.send-btn:disabled {
	background-color: var(--tn-comment);
	color: var(--tn-bg-darker);
	cursor: default;
}
</style>
