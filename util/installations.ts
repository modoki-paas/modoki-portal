export interface User {
  login: string;
  id: number;
  node_id: string;
  avatar_url: string;
  html_url: string;
  gravatar_id: string;
  type: string;
  site_admin: boolean;
  url: string;
  events_url: string;
  following_url: string;
  followers_url: string;
  gists_url: string;
  organizations_url: string;
  received_events_url: string;
  repos_url: string;
  starred_url: string;
  subscriptions_url: string;
}

export interface Permissions {
  checks: string;
  contents: string;
  deployments: string;
  metadata: string;
  pull_requests: string;
}

export interface Installation {
  id: number;
  app_id: number;
  target_id: number;
  account: User;
  access_tokens_url: string;
  repositories_url: string;
  html_url: string;
  target_type: string;
  repository_selection: string;
  events: string[];
  permissions: Permissions;
  created_at: Date;
  updated_at: Date;
}

export interface Repository {
  name: string;
  owner: User;
}

export async function listInstallations(): Promise<Installation[]> {
  return await (await fetch("/github/installations")).json()
}

export async function listRepositories(installationID: number): Promise<Repository[]> {
  return await (await fetch(`/github/repositories?installation_id=${installationID}`)).json()
}
