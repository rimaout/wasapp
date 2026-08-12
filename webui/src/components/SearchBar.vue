<script>
import { ref, watch } from 'vue';

export default {
	props: {
		modelValue: { type: String, default: '' },
	},
	emits: ['update:modelValue'],
	setup(props, { emit }) {
		const query = ref(props.modelValue);
		watch(query, (val) => emit('update:modelValue', val));
		watch(() => props.modelValue, (val) => { query.value = val; });
		return { query };
	},
};
</script>

<template>
	<div class="px-3 pt-3 pb-2">
		<div class="search-wrapper">
			<svg class="feather search-icon" style="width: 16px; height: 16px;">
				<use href="/feather-sprite-v4.29.0.svg#search"/>
			</svg>
			<input type="text" class="form-control form-control-sm search-input"
				placeholder="Search chats..." v-model="query" />
			<button v-if="query" class="search-clear" @click="query = ''">
				<svg class="feather" style="width: 14px; height: 14px; stroke-width: 2.5;">
					<use href="/feather-sprite-v4.29.0.svg#x"/>
				</svg>
			</button>
		</div>
	</div>
</template>

<style scoped>
.search-wrapper {
	position: relative;
	display: flex;
	align-items: center;
}

.search-icon {
	position: absolute;
	left: 10px;
	color: var(--tn-comment);
	pointer-events: none;
	z-index: 1;
}

.search-input {
	padding-left: 32px;
	padding-right: 32px;
	border-radius: 20px;
}

.search-clear {
	position: absolute;
	right: 6px;
	background: none;
	border: none;
	padding: 2px;
	color: var(--tn-fg-dark);
	cursor: pointer;
	line-height: 1;
	display: flex;
	align-items: center;
}

.search-clear:hover {
	color: var(--tn-fg);
}
</style>
