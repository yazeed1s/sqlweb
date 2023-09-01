<script lang="ts">
	import { page } from '$app/stores';
	import go2 from '$lib/images/go2.png';
</script>

<body>
	<div class="wrapper">
		<div class="not_found">
			<div class="image-container">
				<img src={go2} alt="Error 404" class="image" />
			</div>
			<h1>Oops! {$page?.error?.message}</h1>
		</div>
	</div>
</body>

<style>
	.wrapper {
		display: grid;
	}
	* {
		transition: all 0.6s;
		margin: 0;
		padding: 0;
		box-sizing: border-box;
		font-family: 'Poppins', sans-serif;
	}

	body {
		font-family: 'Lato', sans-serif;
		color: #888;
		height: 100vh;
		display: grid;
		/* justify-content: center; */
		align-items: center;
		background: linear-gradient(135deg, #e4e9ec, #8ea1ad);
	}

	.image-container {
		margin-bottom: 20px;
	}

	.image {
		width: 300px;
		height: auto;
	}

	.not_found {
		display: flex;
		flex-direction: column;
		align-items: center;
	}

	.not_found h1 {
		font-size: 50px;
		display: inline-block;
		padding-right: 12px;
		animation: type 0.5s alternate infinite;
	}

	@keyframes type {
		from {
			box-shadow: inset -3px 0px 0px #888;
		}
		to {
			box-shadow: inset -3px 0px 0px transparent;
		}
	}
</style>
