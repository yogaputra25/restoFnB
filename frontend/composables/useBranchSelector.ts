export const useBranchSelector = () => {
  const selectedBranchId = useState<string | null>('selectedBranchId', () => null)
  const branches = useState<Array<{ id: string; name: string }>>('branches', () => [])

  const { $api } = useNuxtApp()

  async function fetchBranches() {
    try {
      const res = await $api('/api/admin/branches/list')
      branches.value = (res as any).data || []
    } catch { /* branches stay empty */ }
  }

  function setBranch(id: string | null) {
    selectedBranchId.value = id
  }

  function getBranchName(branchId: string): string {
    return branches.value.find(b => b.id === branchId)?.name || branchId.slice(0, 8)
  }

  return { selectedBranchId, branches, fetchBranches, setBranch, getBranchName }
}
