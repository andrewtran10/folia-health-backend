<?php

use App\Controllers\API\UserController;
use App\Controllers\API\ReminderController;
use Illuminate\Support\Facades\Route;

/**
 * Use this file to define new API routes under the /api/... path
 * 
 * Here are some example, user related endpoints we have established as an example
 */

// User endpoints
Route::get('/users/{id}', [UserController::class, 'read']);
Route::post('/users', [UserController::class, 'create']);

// Reminder endpoints
Route::post('/reminders', [ReminderController::class, 'create']);
Route::patch('/reminders/{reminder}', [ReminderController::class, 'update']);
Route::delete('/reminders/{reminder}', [ReminderController::class, 'delete']);
Route::get('/reminders', [ReminderController::class, 'read']);