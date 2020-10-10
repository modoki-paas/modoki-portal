import { MutationTree } from 'vuex';

export const state = () => ({
  namespace: "default",
});

export type State = ReturnType<typeof state>;

export default State;

export const mutations: MutationTree<State> = {
  updateNamespace: (state: State, namespace: string) => {
    state.namespace = namespace;
  }
}
