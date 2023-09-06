<script lang="ts">
	import { goto } from '$app/navigation';
	import { fade, slide } from 'svelte/transition';
	import { connectedStore, schemaName, type Table, tableDataStore } from '../store';
	import { endpoints, httpClient } from '../../services/api';

	interface DatabaseConfig {
		host: string;
		port: any;
		user: string;
		password: string;
		databaseType: string;
		database: string;
	}

	let form: DatabaseConfig = {
		host: '',
		port: 0,
		user: '',
		password: '',
		databaseType: '',
		database: ''
	};

	let savedConnections: DatabaseConfig[] = [];
	let connectionMode: string = 'h'; // Default connection mode
	let response: string = '';
	let error: string = '';
	let msg: string = '';
	let isError: boolean = false;
	let isSuccess: boolean = false;
	let show: boolean = false;
	let isFormEmpty: boolean = true;

	const sendRequest = async (): Promise<void> => {
		form.port = parseInt(form.port, 10);
		form.databaseType = form.databaseType.toLocaleLowerCase();
		const res = await httpClient(endpoints.connect, {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify(form)
		});
		const json = await res.json();
		res.ok ? handleSuccessResponse(json) : handleErrorResponse(json);
	};

	const sendSaveRequest = async (): Promise<void> => {
		form.port = parseInt(form.port, 10);
		const res = await httpClient(endpoints.save, {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify(form)
		});
		const json = await res.json();
		res.ok ? handleSuccessResponse(json) : handleErrorResponse(json);
	};

	const sendSavedConnectionRequest = async (): Promise<void> => {
		form.port = parseInt(form.port, 10);
		const res = await httpClient(endpoints.saveConnection, {
			method: 'GET',
			headers: { 'Content-Type': 'application/json' }
		});
		const json = await res.json();
		if (res.ok) {
			show = !show;
			savedConnections = json.data.map((item: any) => ({
				host: item.host,
				port: item.port,
				user: item.user,
				password: item.password,
				databaseType: item.databaseType,
				database: item.database
			}));
			console.log(savedConnections[0]);
		} else {
			handleErrorResponse(json);
		}
	};

	const handleSuccessResponse = (json: any): void => {
		if (json.data) {
			const tables: Table[] = json.data.tables.map((table: any) => ({
				name: table.table_name,
				columns: table.columns
			}));
			// response = json.message;
			isSuccess = true;
			tableDataStore.set({
				tables: tables
			});
			schemaName.set({
				name: json.data.schema
			});
			connectedStore.set(true);
			goto('/db');
		}
		isError = false;
		isSuccess = true;
		msg = json.message;
	};

	const handleErrorResponse = (json: any): void => {
		isSuccess = false;
		isError = true;
		error = json.error;
		msg = json.message;
	};

	const onFormSubmitHandler = async (): Promise<void> => {
		await sendRequest();
	};

	const setConnection = (conn: DatabaseConfig): void => {
		form = conn;
	};

	$: {
		isFormEmpty =
			form.host === '' ||
			form.port === 0 ||
			form.user === '' ||
			form.databaseType === '' ||
			form.database === '';
	}
</script>

<div class="container">
	<div class="content">
		<form method="POST" on:submit|preventDefault={onFormSubmitHandler}>
			<div class="user-details">
				<div class="input-box">
					<span class="details">Database type</span>
					<select class="select" bind:value={form.databaseType}>
						<option value="MySQL">MySQL</option>
						<option value="postgreSQL">PostgreSQL</option>
						<option value="SQLite">SQLite</option>
					</select>
				</div>
				<div class="input-box">
					<span class="details">Connection mode</span>
					<select class="select" bind:value={connectionMode}>
						<option value="h">Host and Port</option>
						<option value="s">Socket</option>
					</select>
				</div>
				{#if connectionMode === 'h'}
					<div class="input-box">
						<span class="details">Host</span>
						<input
							bind:value={form.host}
							class="input"
							name="host"
							type="text"
							placeholder="Enter host..."
						/>
					</div>
					<div class="input-box">
						<span class="details">Port</span>
						<input
							bind:value={form.port}
							class="input"
							name="port"
							type="text"
							placeholder="Enter port..."
						/>
					</div>
				{:else if connectionMode === 's'}
					<div class="input-box">
						<span class="details">Socket path</span>
						<input
							class="input"
							name="socketPath"
							type="text"
							placeholder="Enter socket path..."
						/>
					</div>
				{/if}
				<div class="input-box">
					<span class="details">User</span>
					<input
						bind:value={form.user}
						class="input"
						name="name"
						type="text"
						placeholder="Enter username..."
					/>
				</div>
				<div class="input-box">
					<span class="details">Password</span>
					<input
						bind:value={form.password}
						class="input"
						name="password"
						type="password"
						placeholder="Enter password..."
					/>
				</div>
				<div class="input-box">
					<span class="details">Default database</span>
					<input
						bind:value={form.database}
						class="input"
						name="database"
						type="text"
						placeholder="Enter database name..."
					/>
				</div>
			</div>
		</form>
		{#if isError}
			<div class="error" role="alert">
				<p class="font-error">Error: {error}</p>
				<p class="msg">{msg}</p>
			</div>
		{/if}
		{#if isSuccess}
			<div class="success" role="alert">
				<p class="font">{msg}</p>
			</div>
		{/if}
		<div class="button-container">
			<div class="button">
				<input
					type="submit"
					value="Connect"
					disabled={isFormEmpty}
					on:click={sendRequest}
				/>
			</div>
			<div class="button">
				<input
					type="submit"
					value="Save connection"
					disabled={isFormEmpty}
					on:click={sendSaveRequest}
				/>
			</div>
		</div>
		<div class="saved-button-container">
			<div class="button">
				<input
					style="width: 185px; align-self: flex-start; color: rgb(var(--color-surface-100)); background-color: rgb(var(--color-surface-500));"
					type="submit"
					value="View saved connections"
					on:click={sendSavedConnectionRequest}
				/>
			</div>
		</div>
	</div>
	{#if show}
		<!-- svelte-ignore missing-declaration -->
		<div class="saved-connections-box" transition:fade={{duration: 200}}>
			{#each savedConnections as conn}
				<!-- svelte-ignore a11y-click-events-have-key-events -->
				<div class="single-connection" on:click={() => setConnection(conn)}>
					<span>{conn.databaseType} / {conn.user} / {conn.database}</span>
				</div>
			{/each}
		</div>
	{/if}
</div>

<style>
	.saved-connections-box {
		margin-top: 4rem;
		background-color: rgb(var(--color-surface-800));
		/* width: 30%; */
		height: 24vw;
		overflow-y: scroll;
	}

	.saved-connections-box .single-connection {
		margin-top: 10px;
		margin: 10px;
		height: 30px;
		background-color: rgb(var(--color-surface-900));
		border-left-width: 3px;
		border-color: rgb(var(--color-warning-500));
		display: flex;
		align-items: center;
		justify-content: center;
		cursor: pointer;
	}

	.saved-connections-box .single-connection:hover {
		cursor: pointer;
		border-color: rgb(var(--color-success-500));
		background-color: rgb(var(--color-surface-700));
	}

	.error {
		background-color: rgb(var(--color-surface-600));
		border-left-width: 4px;
		border-color: rgb(var(--color-error-600));
		padding: 0.7rem;
	}

	.error .font-error {
		font-weight: 700;
		color: rgb(var(--color-error-600));
	}

	.error .msg {
		font-weight: 500;
		color: rgb(var(--color-surface-100));
	}

	.success {
		background-color: rgb(var(--color-surface-500));
		border-left-width: 4px;
		border-color: rgb(var(--color-success-500));
		padding: 0.7rem;
	}

	.success .font {
		font-weight: 700;
		color: rgb(var(--color-surface-100));
	}

	* {
		color: rgb(var(--color-surface-100));
	}

	.container {
		background-color: rgb(var(--color-surface-800));
		height: 50vh;
		min-width: 100%;
		padding-left: 1.5vw;
		padding-right: 1.5vw;
	}

	.content {
		display: flex;
		flex-direction: column;
		min-width: 100%;
		padding-top: 2vh;
	}

	.content form .user-details {
		display: flex;
		flex-wrap: wrap;
		justify-content: space-between;
		margin: 20px 0 12px 0;
	}

	form .user-details .input-box {
		margin-bottom: 15px;
		width: calc(100% / 2 - 20px);
	}

	form .input-box span.details {
		display: block;
		font-weight: 500;
		margin-bottom: 5px;
	}

	.user-details .input-box select,
	.user-details .input-box input {
		height: 45px;
		width: 100%;
		outline: none;
		border-radius: 10px;
		padding-left: 15px;
		border: 1px solid rgb(var(--color-surface-500));
		border-bottom-width: 2px;
		transition: all 0.3s ease;
		background-color: rgb(var(--color-surface-100));
	}

	.user-details .input-box select:focus,
	.user-details .input-box input:focus,
	.user-details .input-box select:valid,
	.user-details .input-box input:valid {
		border-color: rgb(var(--color-surface-500));
		background-color: rgb(var(--color-surface-600));
	}

	.user-details .input-box select::placeholder,
	.user-details .input-box input::placeholder {
		color: rgb(var(--color-surface-100));
	}

	select:focus,
	select:valid,
	input[type='password'],
	input[type='text'] {
		background-color: rgb(var(--color-surface-100));
	}

	.content .button-container {
		margin: 2rem;
		display: flex;
		flex-direction: row;
		align-self: center;
	}

	.saved-button-container {
		/* margin: 0; */
		position: absolute;
		bottom: 0;
		margin-bottom: 2rem;
	}

	.button input {
		margin-top: 1vh;
		margin-left: 2vh;
		height: 38px;
		width: 130px;
		border-radius: 10px;
		border: none;
		font-weight: 500;
		letter-spacing: 1px;
		cursor: pointer;
		background: rgb(var(--color-primary-500));
		color: rgb(var(--color-surface-800));
	}

	.button input:disabled {
		cursor: not-allowed;
	}

	.button input:not(:disabled):hover {
		/* background: linear-gradient(135deg, #e4e9ec, #b6c7d1); */
		background: rgb(var(--color-primary-700));
	}

	@media (max-width: 584px) {
		form .user-details .input-box {
			margin-bottom: 15px;
			width: 100%;
		}

		.content form .user-details {
			max-height: 300px;
			overflow-y: scroll;
		}

		.user-details::-webkit-scrollbar {
			width: 5px;
		}
	}
</style>
