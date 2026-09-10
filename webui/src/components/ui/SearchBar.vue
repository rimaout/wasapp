<script>
import { ref, watch } from 'vue';

/**
 * SearchBar — a reusable search input component with a magnifying glass icon and clear button.
 *
 * Props:
 *  - modelValue:  The text value bound via v-model from the parent component.
 *  - placeholder: Placeholder text shown when the input is empty (default: 'Search chats...').
 *  - compact:     Reduces vertical padding around the component when true (default: false).
 *  - light:       Applies a lighter color palette when true (useful for dark backgrounds, default: false).
 *
 * Emits:
 *  - update:modelValue: Emitted to the parent component whenever the input value changes.
 */
export default {
	props: {
		modelValue:  { type: String,  default: ''                },
		placeholder: { type: String,  default: 'Search chats...' },
		compact:     { type: Boolean, default: false             },
		light:       { type: Boolean, default: false             },
	},

	emits: ['update:modelValue'],
	setup(props, { emit }) {
		// Local variable synced with input field
		const query = ref(props.modelValue);

		// Every time 'query' changes (as the user types), emit an event to the parent (update v-model binding)
		watch(query, (newVal) => {
			emit('update:modelValue', newVal);
		});

		// If the parent updates 'modelValue' from the outside, sync local 'query'
		watch(() => props.modelValue, (newVal) => {
			query.value = newVal;
		});

		return { query };
	},
};
</script>

<template>
	<div :class="compact ? 'px-3 pb-2' : 'px-3 pt-3 pb-2'">
		<div class="search-wrapper">

			<!-- Search magnifying glass icon overlay -->
			<svg class="feather search-icon" :class="{ light: light }" style="width: 18px; height: 18px;">
				<use href="/feather-sprite-v4.29.0.svg#search"/>
			</svg>

			<!-- Text input field linked to local query variable -->
			<input type="text" class="form-control form-control-sm search-input" :class="{ light: light }"
				:placeholder="placeholder" v-model="query" />

			<!-- Clear button: only visible when text is typed (query is non-empty) -->
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
	color: var(--theme-comment);
	pointer-events: none;
	z-index: 1;
}

.search-icon.light {
	color: var(--theme-fg-dark);
}

.search-input {
	height: 40px;
	padding-left: 36px;
	padding-right: 36px;
	border-radius: 20px;
	border: none;
	font-size: 0.95rem;
}

.search-input.light {
	border: 1px solid var(--theme-border);
	background-color: var(--theme-overlay);
	color: var(--theme-fg);
}

.search-input.light::placeholder {
	color: var(--theme-fg-dark);
}

.search-input.light:focus {
	background-color: var(--theme-overlay);
	color: var(--theme-fg);
}

.search-clear {
	position: absolute;
	right: 6px;
	background: none;
	border: none;
	padding: 2px;
	color: var(--theme-fg-dark);
	cursor: pointer;
	line-height: 1;
	display: flex;
	align-items: center;
}

.search-clear:hover {
	color: var(--theme-fg);
}
</style>
