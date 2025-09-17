<?php

namespace App\Responses;

use Illuminate\Contracts\Support\Responsable;
use Illuminate\Http\JsonResponse;

class NoContentResponse implements Responsable
{
    public function toResponse($request): JsonResponse
    {
        return response()->json(null, 204);
    }
}