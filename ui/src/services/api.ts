export const endpoints = {
	connect: '/connect',
	disconnect: '/disconnect',
	save: '/save',
	saveConnection: '/save/connection',
	execute: '/execute',
	schemas: '/schemas',
	update: '/update'
};

export const httpClient = async (url: string, options: RequestInit = {}) => {
	const response = await fetch(url, options);
	if (!response.ok) throw new Error(response.statusText);
	return await response.json();
};
