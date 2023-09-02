<script lang="ts">
  import { get } from 'svelte/store';
  import { login, register } from '../authentication';
  import { initiateAuthentication } from '../secure-client';
  import { state } from '../state';
  import AuthorizationProgress from './AuthenticationProgress.svelte';
  import Authentication from './forms/Authentication.svelte';
  import Authorization from './forms/Authorization.svelte';
  import Identification from './forms/Identification.svelte';
  import Registration from './forms/Registration.svelte';

  const onInitiateAuthentication = (userId: string): void => {
      state.update((state) => ({
          ...state,
          loading: true,
      } ));

      initiateAuthentication(userId)
          .then((ops) => {
              if (ops.type === 'create') {
                  state.update((oldState) => ({
                      ...oldState,
                      state: 'createCredential',
                      userId,
                      credentialOptions: ops,
                      loading: false,
                  }));
              } else {
                  state.update((oldState) => ({
                      ...oldState,
                      state: 'getCredential',
                      userId,
                      credentialOptions: ops,
                      loading: false,
                  }));
              }
          })
          .catch((err: Error) => {
              state.update((oldState) => ({
                  ...oldState,
                  error: err,
                  loading: false,
              }));
          });
  };

  const onCreateCredential = (): void => {
      const localState = get(state);
      if (localState.state !== 'createCredential') {
          return;
      }

      register(localState.credentialOptions)
          .then(() => {
              state.update((oldState) => ({
                  ...oldState,
                  state: 'success',
              } as any));
          })
          .catch((err: Error) => {
              state.update((oldState) => ({
                  ...oldState,
                  error: err,
              }));
          });
  };

  const onGetCredential = (): void => {
      const localState = get(state);
      if (localState.state !== 'getCredential') {
          return;
      }

      login(localState.credentialOptions)
          .then(() => {
              state.update((oldState) => ({
                  ...oldState,
                  state: 'success',
              } as any));
          })
          .catch((err: Error) => {
              state.update((oldState) => ({
                  ...oldState,
                  error: err,
              }));
          });
  };

  const onAuthorize = (): void => {
    fetch('/v1/oauth2/authorization/succeed')
  }
</script>

<div class="flex justify-between items-center">
  <h1>
    Login
  </h1>
  {#if $state.state === 'userId'}
  <a
    class="text-yellow-500 dark:text-yellow-400 underline rounded"
    href="#foo"
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
<Authentication userId={$state.userId} submit={onGetCredential} />
{:else if $state.state === 'success'}
<Authorization submit={onAuthorize}></Authorization>
{/if}
