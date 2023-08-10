<script lang="ts">
  import Scope from "./Scope.svelte";
  import type { Scope as ScopeType } from "./models";

  export let scopes: ScopeType[];
</script>

<ul class="space-y-1 pt-2">
  {#each scopes as scope}
  <li>
    <Scope scope={scope} />
  </li>
  {/each}
</ul>