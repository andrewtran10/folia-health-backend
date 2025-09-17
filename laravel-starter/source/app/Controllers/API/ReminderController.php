<?php

namespace App\Controllers\API;

use App\Controllers\Controller;
use App\Models\Reminder;
use App\Models\User;
use App\Resources\ReminderResource;

use Illuminate\Contracts\Support\Responsable;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;

use RRule\RRule;

class ReminderController extends Controller
{
    /**
     * Read an existing reminder from the database
     * 
     */
    public function read(Request $request): Responsable
    {
        // Validate request
        $validated = $request->validate([
            'userId' => 'required|exists:users,id',
            'start_date' => 'required_with:end_date|date',
            'end_date' => 'required_with:start_date|date',
            'search' => 'sometimes|string',
        ]);

        // Fetch user's reminders
        $user = User::findOrFail($validated['userId']);

        $query = $user->reminders();

        // Apply filters
        if (isset($validated['search'])) {
            $query->whereRaw('LOWER(description) LIKE ?', ['%' . strtolower($validated['search']) . '%']);
        }

        if (isset($validated['start_date']) && isset($validated['end_date'])) {
            $query->where('start_at', '<=', $validated['end_date']);
        }   

        $reminders = $query->get();

       return ReminderResource::collection($reminders);
    }

    /**
     * Creates a new reminder & stores it in the database
     * 
     */
    public function create(Request $request): Responsable
    {
        // Validate request body
        $validated = $request->validate([
            'user_id' => 'required|exists:users,id',
            'rrule' => 'required|string',
            'description' => 'required|string',
            'start_at' => 'required|date',
        ]);

        $userId = $validated['user_id'];
        $rrule = $validated['rrule'];
        $description = $validated['description'];
        $startAt = $validated['start_at'];

        return new ReminderResource(Reminder::create([
            'user_id' => $userId,
            'rrule' => $rrule,
            'description' => $description,
            'start_at' => $startAt,
        ]));
    }

    public function update(Reminder $reminder, Request $request): Responsable
    {
        // Validate request body
        $validated = $request->validate([
            'user_id' => 'sometimes|exists:users,id',
            'rrule' => 'sometimes|string',
            'description' => 'sometimes|string',
            'start_at' => 'sometimes|date',
        ]);

        $reminder->update($validated);

        return new ReminderResource($reminder);
    }

    public function delete(Reminder $reminder): JsonResponse
    {
        $reminder->delete();

        return response()->json(null, 204);
    }
}
