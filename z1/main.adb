-- Import the Ada.Text_IO and Ada.Containers.Vectors packages
with Ada.Text_IO; use Ada.Text_IO;
with Ada.Containers.Vectors; use Ada.Containers;

-- Define the Main procedure
procedure Main is
   -- Define a record type for a traveler
   type Position is array(1..2) of Integer;

   type Traveler is record
      ID: Integer;
      Pos: Position;
   end record;

   -- Define a record type for a grid
   type GridMap is array(1..10, 1..10) of Integer;

   type Grid is record
      Width, Height: Integer;
      Travelers: Vector(Positive range <>) of Traveler;
      GridMap: GridMap;
      visitedEdges: Map(array(1..2) of Integer, Vector(1..4) of Boolean);
      mu: System.Task_Primitives.Mutex;
   end record;

   -- Define a procedure to move a traveler
   procedure MoveTraveler(g: access Grid; t: access Traveler) is
      -- Define an array of possible moves
      moves: array(1..4, 1..2) of Integer := ((-1, 0), (1, 0), (0, -1), (0, 1));
      -- Choose a random move direction
      moveDir: Integer := Integer(Rand(1, 4));
      -- Get the move vector for the chosen direction
      move: array(1..2) of Integer := moves(moveDir);
      -- Calculate the new position of the traveler
      newPos: array(1..2) of Integer := (t.Pos(1) + move(1), t.Pos(2) + move(2));
      begin
         g.mu.Lock;
         if newPos(1) >= 1 and newPos(1) <= g.Width and newPos(2) >= 1 and newPos(2) <= g.Height and g.GridMap(newPos(1), newPos(2)) = 0 then
            g.GridMap(t.Pos(1), t.Pos(2)) := 0;
            g.GridMap(newPos(1), newPos(2)) := t.ID;
            t.Pos := newPos;
            temp := g.visitedEdges.Find(t.Pos);
            temp(moveDir) := True;
            g.visitedEdges.Replace(t.Pos, temp);
         end if;
         g.mu.Unlock;
      end MoveTraveler;

   -- Define a procedure to add a traveler to the grid
   procedure AddTraveler(g: access Grid) is
      i: Integer := g.Travelers.Length + 1;
      pos: array(1..2) of Integer := (Integer(Rand(1, g.Width)), Integer(Rand(1, g.Height)));
      t: Traveler := (ID => i, Pos => pos);
   begin
      loop
         g.mu.Lock;
         if g.GridMap(pos(1), pos(2)) = 0 then
            g.GridMap(t.Pos(1), t.Pos(2)) := t.ID;
            g.Travelers.Append(t);
            g.mu.Unlock;
            declare
               task Move_Traveler_Task is
                  pragma Task_Storage_Size(0);
               begin
                  loop
                     delay 1.0;
                     MoveTraveler(g'Access, t'Access);
                  end loop;
               end Move_Traveler_Task;
            begin
               null;
            end;
            exit;
         else
            g.mu.Unlock;
         end if;
      end loop;
   end AddTraveler;

   -- Define a procedure to take a photo of the grid
   procedure TakePhoto(g: access Grid) is
   begin
      g.mu.Lock;
      Put_Line("Grid:");
      Put_Line(" 00 01 02 03 04 05 06 07 08 09");
      for i in 1..2*g.Width-1 loop
         if i mod 2 = 0 then
            Put(i/2, 0);
         else
            Put(" ", 0);
         end if;
         for j in 1..2*g.Height-1 loop
            if i mod 2 /= 0 and j mod 2 /= 0 then
               Put("+", 0);
            elsif i mod 2 /= 0 then
               if g.visitedEdges.Find((i+1)/2, j/2)(2) or g.visitedEdges.Find((i-1)/2, j/2)(1) then
                  Put("==", 0);
               else
                  Put("--", 0);
               end if;
            elsif j mod 2 /= 0 then
               if g.visitedEdges.Find(i/2, (j+1)/2)(4) or g.visitedEdges.Find(i/2, (j-1)/2)(3) then
                  Put(";", 0);
               else
                  Put("|", 0);
               end if;
            else
               if g.GridMap(i/2, j/2) = 0 then
                  Put("  ", 0);
               else
                  Put(g.GridMap(i/2, j/2), 2);
               end if;
            end if;
         end loop;
         New_Line;
      end loop;
      for i in 1..g.Height loop
         for j in 1..g.Width loop
            g.visitedEdges.Replace((i, j), (others => False));
         end loop;
      end loop;
      g.mu.Unlock;
   end TakePhoto;

   -- Create a new grid
   g: aliased Grid := (
      Width => 10,
      Height => 10,
      Travelers => Vector(1..10),
      GridMap => (others => (others => 0)),
      visitedEdges => (others => (others => False)),
      mu => System.Task_Primitives.Mutex
   );

begin
   -- Add 10 travelers to the grid
   for i in 1..10 loop
      AddTraveler(g'Access);
      delay Duration(Rand(0, 10));
      if g.Width * g.Height = g.Travelers.Length then
         exit;
      end if;
   end loop;

   -- Take a photo of the grid
   TakePhoto(g'Access);

   -- Continuously take photos of the grid every 5 seconds
   loop
      delay 5.0;
      TakePhoto(g'Access);
   end loop;
end Main;