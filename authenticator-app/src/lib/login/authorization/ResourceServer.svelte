<script lang="ts">
  import ExpandIcon from "../../utils/ExpandIcon.svelte";
import ScopeList from "./ScopeList.svelte";
  import type { ResourceServer } from "./models";

  let open = true;

  export let resourceServer: ResourceServer;
</script>

<div class="">
  <button
    id={resourceServer.id}
    type="button"
    aria-expanded={open}
    aria-controls={`${resourceServer.id}-scopes`}
    class="flex w-full items-center px-2 py-1 border border-stone-500 rounded-lg -outline-offset-4"
    on:click={() => open = !open}
  >
    <img class="w-6 h-6 mr-2 bg-white rounded-md inline-block" src={resourceServer.src} alt={`resource-server-${resourceServer.id}`}>
    <h3 class="inline-block text-lg">{resourceServer.title}</h3>
    <ExpandIcon class={`ml-auto fill-stone-50 ${open ? '' : 'rotate-180'}`} />
  </button>
  <div class={open ? '' : 'hidden'} id={`${resourceServer.id}-scopes`} aria-labelledby={resourceServer.id}>
    <ScopeList scopes={resourceServer.scopes} />
  </div>
</div>
