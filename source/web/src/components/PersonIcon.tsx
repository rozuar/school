import React from 'react'

export function PersonIcon({
  color = '#16a34a',
  size = 18,
  title,
}: {
  color?: string
  size?: number
  title?: string
}) {
  return (
    <svg
      width={size}
      height={size}
      viewBox="0 0 24 24"
      aria-label={title}
      role="img"
      style={{ display: 'inline-block', verticalAlign: 'text-bottom' }}
    >
      {title ? <title>{title}</title> : null}
      <path
        fill={color}
        d="M12 12c2.76 0 5-2.24 5-5S14.76 2 12 2 7 4.24 7 7s2.24 5 5 5Zm0 2c-4.42 0-8 2.24-8 5v3h16v-3c0-2.76-3.58-5-8-5Z"
      />
    </svg>
  )
}
