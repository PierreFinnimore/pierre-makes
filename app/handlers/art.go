package handlers

import (
	"pierre/app/types"
	"pierre/app/views/art"

	"pierre/kit"
)

func HandleArtIndex(kit *kit.Kit) error {
	poems := []types.PersonalPoem{
		{
			Title: "Between",
			Paragraphs: []string{
				`I will tell her
While my ears still ring
With those arhythmic screeches
That she cannot sing
Her talking voice is cream and peaches
But pitched? I cannot stand the thing
I will tell her`,
				`I won’t tell her
While she sings, she smiles
And that makes me grin
And bear it:
She takes her soul and bears it, while
To shoot the Mockingbird’s a sin
Hers is more the squawking style
I won’t tell her`,
				`I listen and my brain is split
My eyes fixed firmly to the clock
I need to run and yet I sit
Between the hard place and the rock`,
			},
		},
		{
			Title: "Dammed",
			Paragraphs: []string{
				`I am dammed so crops can grow
So they can reap what they can sow
It is rot that feeds the land
Smell the flowers; I am dammed.`,
				`I am dammed by wall of stone
Yet years ago my power alone
Turned stone to pebble, pebble to sand
They built too high, and I am dammed`,
				`Past the world, I used to run:
To feel the wind, to feel the sun.
Before the wall, still I stand
I used to run, but I am dammed`,
				`It is strange, this damming curse,
I do not change, but I feel worse.
I’d climb it, but I’ve no hands.
This wall is mine, and I am dammed.`,
				`What shape was I, before the wall?
Did I flow? Did I fall?
So neat and clean, and all to plan:
I am the wall; I am the dam.`,
			},
		},
	}

	stories := []types.PersonalPoem{
		{
			Title: "Caving",
			Paragraphs: []string{
				`Looking up the rope, tautly umbillicalling me to the darkness of the cave above, I caught myself mentally checking off the safety redundancies for what must've been the third or fourth time since I started the descent: the harness, the knots, the helmet, the rope, the radio, the knife, the proximity alarm, the anchor I'd selected, the backup, and finally Tim; my partner on this mission. As if on cue, his voice, quietened by distance, came from above: 
				"What's the hold up?"
				"Nothing! I was just making sure we'd not misread the depth. I reckon I'm about halfway down now. Everything's ok, no need to worry." I realised I was blabbering now, as I often did when talking to Tim. Something about him set me on edge. Or maybe… "It's fine. I'm fine. We can keep going. I'm starting again now." I descended another 10 or so feet. "Tim?" A pause. Then the rope started shaking, violently. `,
			},
		}}
	return kit.Render(art.Art(poems, stories))
}
