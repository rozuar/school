# Arquitectura del Sistema

## Vision General
Arquitectura modular, orientada a eventos, con separacion clara entre:
- Backend de negocio
- Frontends por rol
- Canales de notificacion
- Capa de datos
- Infraestructura portable (Railway -> GCP)

## Principios de Diseno
- Arquitectura modular y orientada a eventos
- Configuracion por mantenedores (no hardcode)
- Trazabilidad completa (auditoria)
- Escalabilidad cloud-first
- Multi-rol con RBAC estricto
- Portable: demo en Railway, produccion en GCP

## Stack Tecnologico

### Backend - Go (Golang)
- Framework: Gin o Fiber
- ORM: GORM o sqlc
- Auth: JWT + RBAC
- Documentacion: OpenAPI / Swagger

### Frontend
- Profesor / Administrativo: Next.js
- Backoffice: Next.js (app separada)
- Mobile: React Native

### Base de Datos
- PostgreSQL (Railway demo, Cloud SQL en produccion)

### Infraestructura
- Demo: Railway
- Produccion: GCP (Cloud Run, Cloud SQL, Pub/Sub)

## Componentes del Backend
- API Gateway
- Rule Engine (evalua condiciones temporales)
- Event Processor (registra y dispara acciones)
- Notification Service
- Case Management Service

## Infraestructura Demo (Railway)
- Backend Go
- Frontend Next.js (Profesor / Admin)
- Backoffice Next.js
- PostgreSQL
- Redis (opcional para colas)

## Infraestructura Produccion (GCP)
- Cloud Run (Backend Go, Next.js SSR)
- Cloud SQL / AlloyDB
- Pub/Sub (eventos y reglas)
- Cloud Tasks (acciones diferidas)
- Firebase (push notifications)
- Cloud Logging & Monitoring
- Secret Manager
- Cloud Load Balancer
- IAM + Identity Platform

## Seguridad
- RBAC estricto por perfil
- Auditoria completa de eventos
- Cifrado en transito y en reposo
- Separacion de entornos (demo / prod)
- Cumplimiento normativo educacional (datos sensibles)
