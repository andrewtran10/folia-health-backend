ARG PHP_VERSION=8.3.3

FROM public.ecr.aws/docker/library/php:8.3.3-fpm AS base

RUN apt-get update && apt-get install -y \
    g++ \
    zip \
    cron \
    curl \
    unzip \
    libpq-dev \
    libicu-dev \
    zlib1g-dev \
    libonig-dev \
    libxrender1 \
    wkhtmltopdf \
    libmemcached-dev \
    libcurl4-gnutls-dev \
    default-mysql-client

RUN apt-get update && apt-get install -y ca-certificates curl gnupg
RUN curl -fsSL https://deb.nodesource.com/gpgkey/nodesource-repo.gpg.key | gpg --dearmor -o /etc/apt/keyrings/nodesource.gpg
RUN echo "deb [signed-by=/etc/apt/keyrings/nodesource.gpg] https://deb.nodesource.com/node_20.x nodistro main" | tee /etc/apt/sources.list.d/nodesource.list
RUN apt-get update &&  apt-get install nodejs -y

RUN apt-get clean && rm -rf /var/lib/apt/lists/*

RUN docker-php-ext-configure intl \
    && docker-php-ext-install \
    pdo \
    curl \
    intl \
    bcmath \
    mysqli \
    opcache \
    mbstring \
    pdo_mysql

RUN pecl install redis && docker-php-ext-enable redis
RUN pecl install memcached && docker-php-ext-enable memcached
RUN curl -sS https://getcomposer.org/installer | php -- --install-dir=/usr/local/bin --filename=composer

WORKDIR /var/www/html

RUN rm -Rf /var/www/* && mkdir -p /var/www/html

COPY ./bin/php.ini /usr/local/etc/php/
COPY ./bin/opcache.ini /usr/local/etc/php/
COPY ./bin/entrypoint /usr/local/bin/

RUN chmod +x /usr/local/bin/entrypoint

COPY --chown=www-data:www-data . .

RUN mkdir -p ./storage/framework \
    && mkdir -p storage/framework/{cache, testing, sessions, views} \
    && mkdir -p storage/framework/bootstrap \
    && mkdir -p storage/framework/bootstrap/cache

RUN usermod --uid 1000 www-data \
    && groupmod --gid 1000  www-data \
    && chown -R www-data:www-data storage/

FROM base AS app

VOLUME [ "/var/www/html/public" ]

EXPOSE 9000

ENTRYPOINT [ "/usr/local/bin/entrypoint" ]