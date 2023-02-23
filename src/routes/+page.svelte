<script lang="ts">
	import Box from './box.svelte';
	import type { ResponseResult } from './types';

	let targetLink = '/delay/1ms';
	let interval = 100;
	let deactive = true;
	let useOpacity = true;
	let useTransition = true;

	let maxDur = 1.0;
	let minDur = Number.MAX_VALUE;
	let movingAvg = 0;
	let results = Array<ResponseResult>(1000).fill({ isFailed: false });
	let ptr = 0;

	const request = () => {
		if (deactive) return;
		const pptr = ptr;
		ptr += 1;
		ptr %= 1000;

		results[pptr] = { isFailed: false };
		const startedAt = new Date();
		fetch(targetLink, { cache: 'no-store' })
			.then((response) => {
				results[pptr] = {
					isFailed: false,
					startedAt: startedAt,
					duration: new Date().valueOf() - startedAt.valueOf(),
					status: response.status,
					color: response.headers.get('color') || 'aqua'
				};
			})
			.catch((reason) => {
				results[pptr] = {
					isFailed: true,
					startedAt: startedAt,
					duration: new Date().valueOf() - startedAt.valueOf()
				};
				console.log(reason);
			})
			.finally(() => {
				const dur = results[pptr].duration || 0;
				if (maxDur < dur) maxDur = dur;
				if (minDur > dur) minDur = dur;
				movingAvg = movingAvg * 0.7 + dur * 0.3;
				results = results;
			});
	};

	let clear: NodeJS.Timer;
	$: {
		clearInterval(clear);
		clear = setInterval(request, interval);
	}
</script>

<div>
	<h1>Simple Loadtester</h1>
	<input bind:value={targetLink} />
	<button class:deactive on:click={() => (deactive = !deactive)}>Active</button>
	<button on:click={() => (interval = 1)}>1ms</button>
	<button on:click={() => (interval = 10)}>10ms</button>
	<button on:click={() => (interval = 50)}>50ms</button>
	<button on:click={() => (interval = 100)}>100ms</button>
	<button
		on:click={() => {
			results = Array(1000).fill({ isFailed: false });
			maxDur = 1.0;
			minDur = Number.MAX_VALUE;
			movingAvg = 0;
			ptr = 0;
		}}>clear</button
	>
	<button
		on:click={() => {
			maxDur = 1.0;
			minDur = Number.MAX_VALUE;
			movingAvg = 0;
		}}>reset</button
	>
</div>
<div class="group">
	<label>
		<input type="checkbox" bind:checked={useOpacity} />
		visualize delay
	</label>
	<label>
		<input type="checkbox" bind:checked={useTransition} />
		use transition
	</label>
</div>

<div class="group">
	<p>request to: {targetLink}</p>
	<p>interval: {interval}ms</p>
</div>
<div class="group">
	<p>ptr: {ptr}</p>
	<p>moving avg: {movingAvg.toFixed(2)}ms</p>
	<p>max latency: {maxDur}ms</p>
	<p>min latency: {minDur}ms</p>
</div>

<div class="result">
	{#each results as result, i}
		<Box {result} {maxDur} {minDur} selected={ptr === i} {useOpacity} {useTransition} />
	{/each}
</div>

<style>
	.group {
		display: flex;
	}
	.group > p {
		margin-right: 10px;
		width: 170px;
		text-align: left;
		white-space: nowrap;
	}
	.result {
		display: flex;
		flex-direction: row;
		flex-wrap: wrap;
		max-width: 780px;
	}
	.deactive {
		background-color: pink;
	}
</style>
