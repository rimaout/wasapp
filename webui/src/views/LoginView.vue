<script>
import { setAuth } from '../services/auth.js';

export default {
	data() {
		return {
			username: '',
			errormsg: null,
			loading: false,
		}
	},
	methods: {
		async doLogin() {
			this.errormsg = null;
			if (!this.username.trim()) { // Check if the username is empty or only whitespace
				this.errormsg = 'Please enter a username';
				return;
			}
			this.loading = true;
			try {
				// Send a POST request to the /session endpoint with the username
				let response = await this.$axios.post('/session', {
					userName: this.username.trim()
				});

				// Extract the token from the response and remove the 'Bearer ' prefix
				let token = response.data.token.replace('Bearer ', '');
				setAuth(token, response.data.userId);

				// Redirect to the /chats route after successful login
				this.$router.push('/chats');
			} catch (e) {
				this.errormsg = e.response?.data?.message || e.toString();
			}
			this.loading = false;
		},
	},
};

</script>

<template>
	<div class="d-flex align-items-center justify-content-center min-vh-100">
		<div class="login-box">
			<div class="card card-body shadow p-4">
				<div class="text-center mb-4">
					<h1 class="display-5">WASApp</h1>
					<p class="text-muted mb-0">Enter your username to sign in</p>
				</div>

				<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>

				<div class="mb-3">
					<label class="form-label">Username</label>
					<input
						v-model="username"
						type="text"
						class="form-control"
						placeholder="e.g. Feldspar"
						@keyup.enter="doLogin"
						:disabled="loading"
					/>
				</div>

				<button
					class="btn btn-primary w-100"
					@click="doLogin"
					:disabled="loading"
				>
					<span v-if="loading" class="spinner-border spinner-border-sm me-2"></span>
					Sign in
				</button>
			</div>
		</div>
	</div>
</template>


<style scoped>
.login-box {
	width: 400px;
	max-width: 90vw;
}
</style>
