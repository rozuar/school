import type { Metadata } from 'next'
import './globals.css'

export const metadata: Metadata = {
  title: 'Monitoreo de Salas - Dashboard Inspectoría',
  description: 'Sistema de monitoreo centralizado de salas de clases',
}

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="es">
      <body>{children}</body>
    </html>
  )
}
