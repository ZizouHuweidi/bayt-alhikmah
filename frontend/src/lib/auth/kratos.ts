import type {
	LoginFlow,
	RegistrationFlow,
	RecoveryFlow,
	VerificationFlow,
	SettingsFlow,
	Session,
	Identity,
	UiNode,
} from '@ory/client'

const KRATOS_PUBLIC_URL = import.meta.env.VITE_KRATOS_PUBLIC_URL || 'http://localhost:4433'

interface FlowResponse<T> {
	flow: T
	status: number
}

interface ErrorResponse {
	error: string
	message?: string
}

// Helper to get CSRF token from flow
export function getCsrfToken(nodes: UiNode[]): string {
	const csrfNode = nodes.find(
		(node) =>
			node.attributes &&
			'node_type' in node.attributes &&
			node.attributes.node_type === 'input' &&
			'name' in node.attributes &&
			node.attributes.name === 'csrf_token',
	)

	if (
		csrfNode &&
		csrfNode.attributes &&
		'value' in csrfNode.attributes &&
		csrfNode.attributes.value
	) {
		return String(csrfNode.attributes.value)
	}

	return ''
}

// Create a new login flow
export async function createLoginFlow(
	returnTo?: string,
): Promise<FlowResponse<LoginFlow> | ErrorResponse> {
	const params = new URLSearchParams()
	if (returnTo) params.append('return_to', returnTo)

	const response = await fetch(
		`${KRATOS_PUBLIC_URL}/self-service/login/browser?${params.toString()}`,
		{
			method: 'GET',
			credentials: 'include',
			headers: {
				Accept: 'application/json',
			},
		},
	)

	if (!response.ok) {
		return {
			error: 'Failed to create login flow',
			message: await response.text(),
		}
	}

	return {
		flow: await response.json(),
		status: response.status,
	}
}

// Get existing login flow
export async function getLoginFlow(flowId: string): Promise<FlowResponse<LoginFlow> | ErrorResponse> {
	const response = await fetch(
		`${KRATOS_PUBLIC_URL}/self-service/login/flows?id=${flowId}`,
		{
			method: 'GET',
			credentials: 'include',
			headers: {
				Accept: 'application/json',
			},
		},
	)

	if (!response.ok) {
		return {
			error: 'Failed to get login flow',
			message: await response.text(),
		}
	}

	return {
		flow: await response.json(),
		status: response.status,
	}
}

// Submit login form
export async function submitLoginFlow(
	flowId: string,
	data: Record<string, string>,
): Promise<FlowResponse<LoginFlow> | { redirect_to: string } | ErrorResponse> {
	const response = await fetch(
		`${KRATOS_PUBLIC_URL}/self-service/login?flow=${flowId}`,
		{
			method: 'POST',
			credentials: 'include',
			headers: {
				'Content-Type': 'application/json',
				Accept: 'application/json',
			},
			body: JSON.stringify(data),
		},
	)

	if (response.status === 422) {
		// Flow requires additional steps
		return {
			flow: await response.json(),
			status: response.status,
		}
	}

	if (response.status === 200 || response.status === 201) {
		// Check if it's a redirect
		const contentType = response.headers.get('content-type')
		if (contentType && contentType.includes('application/json')) {
			const data = await response.json()
			if (data.redirect_to) {
				return { redirect_to: data.redirect_to }
			}
		}

		// Successful login
		return {
			flow: await response.json(),
			status: response.status,
		}
	}

	return {
		error: 'Failed to submit login',
		message: await response.text(),
	}
}

// Create a new registration flow
export async function createRegistrationFlow(
	returnTo?: string,
): Promise<FlowResponse<RegistrationFlow> | ErrorResponse> {
	const params = new URLSearchParams()
	if (returnTo) params.append('return_to', returnTo)

	const response = await fetch(
		`${KRATOS_PUBLIC_URL}/self-service/registration/browser?${params.toString()}`,
		{
			method: 'GET',
			credentials: 'include',
			headers: {
				Accept: 'application/json',
			},
		},
	)

	if (!response.ok) {
		return {
			error: 'Failed to create registration flow',
			message: await response.text(),
		}
	}

	return {
		flow: await response.json(),
		status: response.status,
	}
}

// Get existing registration flow
export async function getRegistrationFlow(
	flowId: string,
): Promise<FlowResponse<RegistrationFlow> | ErrorResponse> {
	const response = await fetch(
		`${KRATOS_PUBLIC_URL}/self-service/registration/flows?id=${flowId}`,
		{
			method: 'GET',
			credentials: 'include',
			headers: {
				Accept: 'application/json',
			},
		},
	)

	if (!response.ok) {
		return {
			error: 'Failed to get registration flow',
			message: await response.text(),
		}
	}

	return {
		flow: await response.json(),
		status: response.status,
	}
}

// Submit registration form
export async function submitRegistrationFlow(
	flowId: string,
	data: Record<string, string>,
): Promise<FlowResponse<RegistrationFlow> | { redirect_to: string } | ErrorResponse> {
	const response = await fetch(
		`${KRATOS_PUBLIC_URL}/self-service/registration?flow=${flowId}`,
		{
			method: 'POST',
			credentials: 'include',
			headers: {
				'Content-Type': 'application/json',
				Accept: 'application/json',
			},
			body: JSON.stringify(data),
		},
	)

	if (response.status === 422) {
		return {
			flow: await response.json(),
			status: response.status,
		}
	}

	if (response.status === 200 || response.status === 201) {
		const contentType = response.headers.get('content-type')
		if (contentType && contentType.includes('application/json')) {
			const data = await response.json()
			if (data.redirect_to) {
				return { redirect_to: data.redirect_to }
			}
		}

		return {
			flow: await response.json(),
			status: response.status,
		}
	}

	return {
		error: 'Failed to submit registration',
		message: await response.text(),
	}
}

// Create recovery flow
export async function createRecoveryFlow(): Promise<
	FlowResponse<RecoveryFlow> | ErrorResponse
> {
	const response = await fetch(
		`${KRATOS_PUBLIC_URL}/self-service/recovery/browser`,
		{
			method: 'GET',
			credentials: 'include',
			headers: {
				Accept: 'application/json',
			},
		},
	)

	if (!response.ok) {
		return {
			error: 'Failed to create recovery flow',
			message: await response.text(),
		}
	}

	return {
		flow: await response.json(),
		status: response.status,
	}
}

// Submit recovery form
export async function submitRecoveryFlow(
	flowId: string,
	data: Record<string, string>,
): Promise<FlowResponse<RecoveryFlow> | ErrorResponse> {
	const response = await fetch(
		`${KRATOS_PUBLIC_URL}/self-service/recovery?flow=${flowId}`,
		{
			method: 'POST',
			credentials: 'include',
			headers: {
				'Content-Type': 'application/json',
				Accept: 'application/json',
			},
			body: JSON.stringify(data),
		},
	)

	if (!response.ok) {
		return {
			error: 'Failed to submit recovery',
			message: await response.text(),
		}
	}

	return {
		flow: await response.json(),
		status: response.status,
	}
}

// Create verification flow
export async function createVerificationFlow(): Promise<
	FlowResponse<VerificationFlow> | ErrorResponse
> {
	const response = await fetch(
		`${KRATOS_PUBLIC_URL}/self-service/verification/browser`,
		{
			method: 'GET',
			credentials: 'include',
			headers: {
				Accept: 'application/json',
			},
		},
	)

	if (!response.ok) {
		return {
			error: 'Failed to create verification flow',
			message: await response.text(),
		}
	}

	return {
		flow: await response.json(),
		status: response.status,
	}
}

// Submit verification form
export async function submitVerificationFlow(
	flowId: string,
	data: Record<string, string>,
): Promise<FlowResponse<VerificationFlow> | ErrorResponse> {
	const response = await fetch(
		`${KRATOS_PUBLIC_URL}/self-service/verification?flow=${flowId}`,
		{
			method: 'POST',
			credentials: 'include',
			headers: {
				'Content-Type': 'application/json',
				Accept: 'application/json',
			},
			body: JSON.stringify(data),
		},
	)

	if (!response.ok) {
		return {
			error: 'Failed to submit verification',
			message: await response.text(),
		}
	}

	return {
		flow: await response.json(),
		status: response.status,
	}
}

// Get current session
export async function getSession(): Promise<Session | null> {
	try {
		const response = await fetch(`${KRATOS_PUBLIC_URL}/sessions/whoami`, {
			method: 'GET',
			credentials: 'include',
			headers: {
				Accept: 'application/json',
			},
		})

		if (!response.ok) {
			return null
		}

		return await response.json()
	} catch {
		return null
	}
}

// Logout
export async function logout(): Promise<boolean> {
	try {
		// First, get the logout URL
		const response = await fetch(
			`${KRATOS_PUBLIC_URL}/self-service/logout/browser`,
			{
				method: 'GET',
				credentials: 'include',
				headers: {
					Accept: 'application/json',
				},
			},
		)

		if (!response.ok) {
			return false
		}

		const data = await response.json()

		// Perform the logout
		const logoutResponse = await fetch(data.logout_url, {
			method: 'GET',
			credentials: 'include',
		})

		return logoutResponse.ok
	} catch {
		return false
	}
}

// Check if user is authenticated
export async function isAuthenticated(): Promise<boolean> {
	const session = await getSession()
	return session !== null && session.active === true
}

// Get current identity
export async function getIdentity(): Promise<Identity | null> {
	const session = await getSession()
	return session?.identity || null
}

// Get user traits
export function getUserTraits(identity: Identity | null): {
	email?: string
	firstName?: string
	lastName?: string
} {
	if (!identity || !identity.traits) {
		return {}
	}

	const traits = identity.traits as {
		email?: string
		name?: { first?: string; last?: string }
	}

	return {
		email: traits.email,
		firstName: traits.name?.first,
		lastName: traits.name?.last,
	}
}
