import { type Writable, writable } from 'svelte/store';
import type TableComponent from '$lib/components/Table.svelte';

interface Column {
	field: string;
	type: string;
	key: string;
	constraint_name: string;
	refrenced_table: string;
	refrenced_column: string;
}

export interface Table {
	name: string;
	columns: Column[];
}

export interface Tab {
	name: string;
	component: TableComponent;
}


export const tabStore: Writable<{ tabs: Tab[] }> = writable({
	tabs: [] as Tab[],
});

export const tableDataStore: Writable<{ tables: Table[] }> = writable({
	tables: [] as Table[],
});

export const schemaName: Writable<{ name: string }> = writable({
	name: ''
});

export const connectedStore: Writable<boolean> = writable(false);

export const selectedTable: Writable<{ tableName: string }> = writable({
	tableName: ''
});