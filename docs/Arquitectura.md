# Decisiones y arquitectura de la aplicación

Architectural style: client/server (3 tiers: presentacion, business logic, db)
Lenguaje para la API: Go
Uso de Docker y Docker Compose
Base de datos relacional: PostgreSQL

Estrategia elegida para el modelo de la db: mapper: sqlc
Notamos la dependencia que tiene el proyecto a esta herramienta, al hacer las consultas utilizando queries. No escala muy bien.

Generacion y apicacion de migraciones para la db: Atlas
Recompilacion y reinicio automatico de la app: Air
Automatizacion en el flujo de desarrollo: Makefile

Herramientas empleadas para la capa de presentación: Templ

Pensamos en diseñar una aplicación multi-contenedor, por un lado la lógica del servidor go y por otra parte la base de datos postgres, 
así obtenemos una ejecución más controlada y aislada en la aplicación.

El contenedor app está construido de tal forma que, mientras esté corriendo, cada vez que se realiza una modificación en el código, ésta se aplique inmediatamente a través de air, agilizando el flujo de desarrollo.

Decidimos incluir únicamente 'air' y 'atlas' dentro del contenedor y tener 'sqlc' y 'templ' en el host local para reducir el tamaño y el tiempo de compilación inicial de la imagen.

Separamos el flujo de docker en desarrollo y producción para tener una mejor organización y facilitar el testing y debugging.

Empleamos un archivo example.env como ejemplo para testear la aplicación sin exponer datos sensibles.