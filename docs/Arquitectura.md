# Decisiones y arquitectura de la aplicación

Pensamos en diseñar una aplicación multi-contenedor, por un lado la lógica del servidor go y por otra parte la base de datos postgres, 
así obtenemos una ejecución más controlada y aislada entre las partes.

El contenedor app está construido de tal forma que, mientras esté corriendo, cada vez que se realiza una modificación en el código, ésta se aplique inmediatamente, agilizando el flujo de desarrollo.
Decidimos incluir únicamente 'air' y 'atlas' dentro del contenedor y tener 'sqlc' y 'templ' en el host local para reducir el tamaño y el tiempo de compilación de la imagen.

Separamos el flujo de docker en desarrollo y producción para tener una mejor organización y facilitar el testing.
