import { create } from 'apisauce';
import type { ApisauceInstance, ApiResponse } from 'apisauce';
import { getGeneralApiProblem, type GeneralApiProblem } from './apiProblem';

/**
 * API Sauce 🔥
 * An API abstraction layer, all methods and errors should be defined here
 */

class Api {
	private apisauce: ApisauceInstance;
	constructor() {
		this.apisauce = create({
			baseURL: 'http://localhost:3000',
			headers: { Accept: 'application/json' }
		});
	}

	public async connect(data: {
		host: string;
		port: number;
		user: string;
		password: string;
		databaseType: string;
		database: string;
	}): Promise<{ data: any } | GeneralApiProblem> {
		const response: ApiResponse<any> = await this.apisauce.post('/connect', data);
		if (!response.ok) {
			const problem = getGeneralApiProblem(response);
			if (problem) return problem;
		}
		return response.data;
	}

	public async disconnect() {
		return this.apisauce.get('/disconnect');
	}

	public async save() {
		return this.apisauce.get('/save');
	}

	public async savedConnection() {
		return this.apisauce.get('/saved/connections');
	}

	public async execute() {
		return this.apisauce.get('/execute');
	}

	public async schemas() {
		return this.apisauce.get('/schemas');
	}

	public async update() {
		return this.apisauce.get('/update');
	}

	public async table() {
		return this.apisauce.get('/table');
	}

	public async export() {
		return this.apisauce.get('/export');
	}
}

export const api = new Api();
export type APIType = typeof api;
