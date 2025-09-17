<?php

namespace App\Controllers\API;

use App\Controllers\Controller;
use App\Models\Reminder;
use App\Models\User;
use App\Resources\ReminderResource;
use App\Requests\ReminderReadRequest;
use App\Requests\ReminderCreateRequest;
use App\Requests\ReminderUpdateRequest;

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
    public function read(ReminderReadRequest $request): Responsable
    {
        // Validate request
        $validated = $request->validated();

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
    public function create(ReminderCreateRequest $request): Responsable
    {
        // Validate request body
        $validated = $request->validated();

        return new ReminderResource(Reminder::create($validated));
    }

    public function update(Reminder $reminder, ReminderUpdateRequest $request): Responsable
    {
        // Validate request body
        $validated = $request->validated();

        $reminder->update($validated);

        return new ReminderResource($reminder);
    }

    public function delete(Reminder $reminder): JsonResponse
    {
        $reminder->delete();

        return response()->json(null, 204);
    }
}
