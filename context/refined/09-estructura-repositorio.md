# Estructura del Repositorio

Basado en `context/raw/carpetas.txt`.

## Carpetas principales

```
school/
├── context/
│   ├── raw/        # Documentos originales (insumos)
│   └── refined/    # Documentación ordenada (salida)
└── source/
    ├── api/        # API Go (negocio + reglas + eventos + websockets)
    ├── web/        # Next.js (Profesor)
    ├── backoffice/ # Next.js (Inspectoría/Admin/Configuración)
    └── android/    # App móvil (base)
```

## Regla práctica

- **`raw/`**: no se edita (sirve como “fuente”).
- **`refined/`**: se edita y organiza (sirve como “manual” del sistema).
