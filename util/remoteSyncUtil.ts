import { DevTsuzuModokiV1alpha1RemoteSyncSpecBase, DevTsuzuModokiV1alpha1RemoteSyncSpec} from "@modoki-paas/kubernetes-fetch-client";

export function calcBase(b?: DevTsuzuModokiV1alpha1RemoteSyncSpecBase): string {
  if(!b) {
    return ""
  }
  let ref = b.github.sha;

  if (!ref && b.github.pullRequest) {
    ref = `#${b.github.pullRequest}`
  }
  if (!ref) {
    ref = b.github.branch ?? "master"
  }

  return `${b.github.owner}/${b.github.repo}(${ref})`
}

export function openApp(spec: DevTsuzuModokiV1alpha1RemoteSyncSpec) {
  const gh = spec.base.github;
  if (gh.sha) {
    window.open(`http://github.com/${gh.owner}/${gh.repo}/tree/${gh.sha}`, '_blank');
  } else if(gh.pullRequest) {
    window.open(`http://github.com/${gh.owner}/${gh.repo}/pull/${gh.pullRequest}`, '_blank');
  } else {
    window.open(`http://github.com/${gh.owner}/${gh.repo}/tree/${gh.branch ?? "master"}`, '_blank');
  }
}
