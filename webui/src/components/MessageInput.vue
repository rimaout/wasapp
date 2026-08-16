<script>
import IconButton from './IconButton.vue';

export default {
	components: { IconButton },
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
			<IconButton class="me-3" icon="image" label="Attach image" />
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
	position: absolute;
	left: 0;
	right: 0;
	bottom: 0;
	padding: 0 24px 14px;
}

.message-input-group {
	display: flex;
	align-items: center;
	background: var(--tn-bg-highlight);
	border-radius: 999px;
	height: 52px;
	padding: 6px 6px 6px 8px;
	box-shadow: 0 3px 16px rgba(0, 0, 0, 0.3);
}

.message-field {
	flex: 1;
	background: none;
	border: none;
	outline: none;
	color: var(--tn-fg);
	font-size: 1rem;
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
}

.send-btn:hover {
	filter: brightness(1.1);
}

.send-btn:disabled {
	background-color: var(--tn-comment);
	color: var(--tn-bg-darker);
	cursor: default;
}
</style>
