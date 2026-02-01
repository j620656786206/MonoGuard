/**
 * Generate an ASCII diagram for a cycle visualization
 * AC2: ASCII diagram for text/markdown exports
 */
export function generateCycleAscii(
  packages: string[],
  breakingPoint: { fromPackage: string; toPackage: string }
): string {
  const shortNames = packages.map((p) => p.split('/').pop() || p)
  const maxLength = Math.max(...shortNames.map((n) => n.length))

  let diagram = '```\n'
  diagram += 'Cycle Path:\n\n'

  for (let i = 0; i < packages.length; i++) {
    const current = shortNames[i]
    const nextIndex = (i + 1) % packages.length
    const next = shortNames[nextIndex]
    const fullCurrent = packages[i]
    const fullNext = packages[nextIndex]

    const isBreaking =
      fullCurrent === breakingPoint.fromPackage && fullNext === breakingPoint.toPackage

    const arrow = isBreaking ? ' ==X==> ' : ' ------> '
    const label = isBreaking ? ' [BREAK HERE]' : ''

    diagram += `  ${current.padEnd(maxLength)}${arrow}${next}${label}\n`

    if (i < packages.length - 1) {
      diagram += `  ${''.padEnd(maxLength)}   |\n`
      diagram += `  ${''.padEnd(maxLength)}   v\n`
    }
  }

  diagram += '\n  (cycle repeats)\n'
  diagram += '```\n'

  return diagram
}
