import { computed } from "vue"
import { useAuthStore } from "~/core/modules/auth/store/AuthStore"

export const useAuth = () => {
	const authStore = useAuthStore()

	const initAuth = async () => {
		authStore.hydrate()
		if (authStore.accessToken) {
			await authStore.refreshProfile()
		}
	}

	return {
		token: computed(() => authStore.accessToken),
		accessToken: computed(() => authStore.accessToken),
		user: computed(() => authStore.currentUser),
		currentUser: computed(() => authStore.currentUser),
		isLoggedIn: computed(() => authStore.isLoggedIn),
		isAuthenticated: computed(() => authStore.isLoggedIn),
		role: computed(() => authStore.role),
		isAdmin: computed(() => authStore.role === "admin"),
		isStaff: computed(() => authStore.role === "staff" || authStore.role === "officer"),
		isFinance: computed(() => authStore.role === "finance"),
		isTenant: computed(() => authStore.role === "tenant"),
		initAuth,
		login: authStore.login,
		logout: authStore.logout,
	}
}
