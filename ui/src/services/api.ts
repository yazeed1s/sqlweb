export const endpoints = {
	connect: '/connect',
	disconnect: '/disconnect',
	save: '/save',
	savedConnection: '/saved/connections',
	execute: '/execute',
	schemas: '/schemas',
	update: '/update',
	table: '/table',
	export: '/export'
};

export const httpClient = async (url: string, options: RequestInit = {}): Promise<Response> => {
	const response = await fetch(url, options);
	// if (!response.ok) throw new Error(response.statusText);
	return response;
};
