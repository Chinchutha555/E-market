export interface IApplication {
	id: number
	application_no?: string
	application_number?: string
	application_type?: "stall_lease" | "general" | string
	subject?: string
	booking_id?: number
	stall_id?: number
	stall_ident: string
	stall_number?: string
	market_id?: number
	market_name?: string
	applicant_name: string
	tenant_name?: string
	applicant_phone: string
	tenant_phone?: string
	applicant_id_card?: string
	applicant_email?: string
	applicant_address?: string
	shop_name: string
	product_category: string
	product_type?: string
	lease_term_months: number
	requested_period_months?: number
	monthly_rate?: number
	stall_monthly_rate?: number
	proposed_rent_type?: string
	requested_start_date?: string
	details?: string
	status: string
	officer_comment?: string
	market_director_comment?: string
	bureau_director_comment?: string
	submitted_at?: string
	reviewed_at?: string
	review_note?: string
	created_at?: string
}

export interface ICreateApplicationPayload {
	application_type?: "stall_lease" | "general" | string
	subject?: string
	booking_id?: number | null
	stall_id?: number
	stall_ident?: string
	stall_number?: string
	market_id?: number
	market_name?: string
	applicant_name: string
	applicant_phone: string
	applicant_id_card?: string
	applicant_email?: string
	applicant_address?: string
	shop_name: string
	product_category?: string
	product_type?: string
	lease_term_months?: number
	requested_period_months?: number
	monthly_rate?: number
	proposed_rent_type?: string
	requested_start_date?: string
	details?: string
	review_note?: string
}
