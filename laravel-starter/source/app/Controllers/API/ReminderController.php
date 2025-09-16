<?php

namespace App\Controllers\API;

use App\Controllers\Controller;
use App\Models\Reminder;
use App\Resources\ReminderResource;
use Illuminate\Contracts\Support\Responsable;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;

class ReminderController extends Controller
{
    /**
     * Read an existing reminder from the database
     * 
     */
    public function read(string $id): Responsable
    {
        return new ReminderResource(Reminder::findOrFail($id));
    }

    /**
     * Creates a new reminder & stores it in the database
     * 
     */
    public function create(Request $request): Responsable
    {
        $userId = $request->input('user_id');
        $rrule = $request->input('rrule');
        $description = $request->input('description');
        $startAt = $request->input('start_at');

        return new ReminderResource(Reminder::create([
            'user_id' => $userId,
            'rrule' => $rrule,
            'description' => $description,
            'start_at' => $startAt,
        ]));
    }
}
