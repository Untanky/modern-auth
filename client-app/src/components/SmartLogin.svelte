<script lang="ts">
  import { get } from 'svelte/store';
  import { initiateAuthentication, signIn, signUp } from '../lib/authentication';
  import { state } from '../lib/state';
  import AuthorizationProgress from './AuthenticationProgress.svelte';
  import Authentication from './forms/Authentication.svelte';
  import Identification from './forms/Identification.svelte';
  import Registration from './forms/Registration.svelte';

  const onInitiateAuthentication = (userId: string): void => {
    state.update(state => ({ ...state, loading: true } ));

    initiateAuthentication(userId)
      .then((ops) => {
        if (ops.type === 'create') {
          state.update(oldState => ({ ...oldState, state: 'createCredential', userId, credentialOptions: ops, loading: false }));
        } else {
          state.update(oldState => ({ ...oldState, state: 'getCredential', userId, credentialOptions: ops, loading: false }));
        }
      })
      .catch((err) => {
        state.update(oldState => ({ ...oldState, error: err, loading: false }));
      });
  };

  const onCreateCredential = (): void => {
    const localState = get(state);
    if (localState.state !== 'createCredential') {
      return;
    }

    signUp(localState.credentialOptions).then(() => {
      state.update(oldState => ({ ...oldState, state: 'success' } as any));
    }).catch((err) => {
      state.update(oldState => ({ ...oldState, error: err }));
    });
  };

  const onGetCredential = (): void => {
    const localState = get(state);
    if (localState.state !== 'getCredential') {
      return;
    }

    signIn(localState.credentialOptions).then(() => {
      state.update(oldState => ({ ...oldState, state: 'success' } as any));
    }).catch((err) => {
      state.update(oldState => ({ ...oldState, error: err }));
    });
  };
</script>

<div class="flex justify-between items-center">
  <h1 class="text-stone-950 dark:text-stone-50 text-2xl font-medium">
    Login
  </h1>
  {#if $state.state === 'userId'}
  <a
    class="text-yellow-500 dark:text-yellow-400 underline"
    href="#"
  >
    Register instead
  </a>
  {/if}
</div>
<AuthorizationProgress></AuthorizationProgress>
{#if $state.state === 'userId'}
<Identification submit={onInitiateAuthentication} />
{:else if $state.state === 'createCredential'}
<Registration submit={onCreateCredential} />
{:else if $state.state === 'getCredential'}
<Authentication submit={onGetCredential} />
{:else if $state.state === 'success'}
<h2 class="text-xl">
  Success
  <p>
    When you click on authenticate a system dialog will open and ask you to authenticate with your biometric data or a physical hardware token. Please prepare for the method chosen when setting up this device.
  </p>
</h2>
{/if}