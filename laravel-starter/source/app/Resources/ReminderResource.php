<?php

namespace App\Resources;

use RRule\RRule;

use Illuminate\Http\Request;
use Illuminate\Http\Resources\Json\JsonResource;

class ReminderResource extends JsonResource
{
    private $rruleObject = null;

    // Initialize RRule Object on resource creation
    public function __construct($resource)
    {
        parent::__construct($resource);
        $this->rruleObject = new RRule($this->resource->rrule);
        // $this->rruleObject = $this->createRRuleObject($this->resource->rrule, $this->resource->start_at);
    }

    public function toArray(Request $request): array
    {
        $data = [
            'id' => $this->id,
            'rrule' => $this->rrule,
            'rrule_human' => $this->transformRRuleToHumanReadable(),
            'description' => $this->description,
            'start_at' => $this->start_at,
        ];

        // If start_date and end_date are provided, include occurrences
        if ($request->has(["start_date", "end_date"])) {
            $start = $request->query('start_date');
            $end = $request->query('end_date');
            $data['occurrences'] = $this->getOccurrencesBetween($start, $end);
        }

        return $data;
    }

    private function getOccurrencesBetween(string $startDate, string $endDate): array
    {
        $start = new \DateTime($startDate);
        $end = new \DateTime($endDate);

        $dates = [];        
        foreach ($this->rruleObject as $occurrence) {           
            if ($occurrence < $start) {
                continue;
            }
            
            if ($occurrence > $end) {
                break;
            }
            
            $dates[] = $occurrence->format('Y-m-d H:i:s');
        }
        return $dates;
    }

    private function transformRRuleToHumanReadable(): string
    {
        return ucwords($this->rruleObject->humanReadable());
    }

    private function createRRuleObject(string $rrule, string $startAt): RRule
    {
        // Remove RRULE: prefix
        $rrule = str_replace('RRULE:', '', $rrule);
        
        // Iterate over each component and build array for RRule Object
        $components = [];
        foreach (explode(';', $rrule) as $part) {
            if (strpos($part, '=') !== false) {
                list($key, $value) = explode('=', $part, 2);
                
                switch ($key) {
                    case 'BYDAY':
                        $components[$key] = explode(',', $value);
                        break;
                    case 'FREQ':
                    case 'WKST':
                        $components[$key] = $value;
                        break;
                    case 'INTERVAL':
                    case 'COUNT':
                    case 'BYHOUR':
                    case 'BYMINUTE':
                    case 'BYMONTHDAY':
                        $components[$key] = (int)$value;
                        break;
                    case 'UNTIL':
                        $components[$key] = new \DateTime($value);
                        break;
                    default:
                        $components[$key] = $value;
                }
            }
        }
        
        $components['DTSTART'] = new \DateTime($startAt);
        
        return new RRule($components);
    }
}
