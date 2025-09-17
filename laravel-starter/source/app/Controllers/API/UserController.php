<?php

namespace App\Controllers\API;

use App\Controllers\Controller;
use App\Models\User;
use App\Resources\UserResource;
use Illuminate\Contracts\Support\Responsable;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;

class UserController extends Controller
{
    /**
     * Read an existing user from the database
     * 
     */
    public function read(string $id): Responsable
    {
        return new UserResource(User::findOrFail($id));
    }

    /**
     * Creates a new user & stores it in the database
     * 
     */
    public function create(Request $request): JsonResponse
    {
        $name = $request->input('name');
        $email = $request->input('email');
        $password = $request->input('password');

        $user = User::create([
            'name' => $name,
            'email' => $email,
            'password' => $password,
        ]);

        $token = $user->createToken('api-key')->plainTextToken;

        return new JsonResponse([
            'user' => new UserResource($user),
            'token' => $token,
        ], 201);
    }
}
