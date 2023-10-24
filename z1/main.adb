with Ada.Text_IO; use Ada.Text_IO;
with Ada.Numerics.Float_Random; use Ada.Numerics.Float_Random;

procedure Main is
    type Position is record
        X : Integer;
        Y : Integer;
    end record;

    type Traveler is record
        ID : Integer;
        Pos : Position;
    end record;

    type Grid is array (1 .. 10, 1 .. 10) of Integer;

    Travelers : array (1 .. 5) of Traveler;
    G : Grid := (others => (others => 0));

    Gen : Ada.Numerics.Float_Random.Generator;

    procedure Add_Traveler(T : Traveler) is
    begin
        if G(T.Pos.X, T.Pos.Y) = 0 then
            G(T.Pos.X, T.Pos.Y) := T.ID;
        else
            Put_Line("Position is already occupied");
        end if;
    end Add_Traveler;

    procedure Move_Traveler(T : in out Traveler) is
        DX : constant array (1 .. 4) of Integer := (1, -1, 0, 0);
        DY : constant array (1 .. 4) of Integer := (0, 0, 1, -1);
        Dir : Integer;
        New_Pos : Position;
    begin
        Dir := Integer(Random(Gen) * 3.0) + 1;
        if T.Pos.X + DX(Dir) in G'Range(1) and then T.Pos.Y + DY(Dir) in G'Range(2) and then G(T.Pos.X + DX(Dir), T.Pos.Y + DY(Dir)) = 0 then
            New_Pos := (X => T.Pos.X + DX(Dir), Y => T.Pos.Y + DY(Dir));
            G(T.Pos.X, T.Pos.Y) := 0;
            G(New_Pos.X, New_Pos.Y) := T.ID;
            T.Pos := New_Pos;
        end if;
    end Move_Traveler;

    procedure Take_Photo is
    begin
        for I in G'Range(1) loop
            for J in G'Range(2) loop
                Put(Integer'Image(G(I, J)) & " ");
            end loop;
            New_Line;
        end loop;
    end Take_Photo;

    protected type Task_Counter is
   procedure Increment;
   procedure Decrement;
   entry Wait_Until_Zero;
private
   Count : Natural := 0;
end Task_Counter;

protected body Task_Counter is
   procedure Increment is
   begin
      Count := Count + 1;
   end Increment;

   procedure Decrement is
   begin
      Count := Count - 1;
      if Count = 0 then
         Wait_Until_Zero'Requeue;
      end if;
   end Decrement;

   entry Wait_Until_Zero when Count = 0 is
   begin
      null;
   end Wait_Until_Zero;
end Task_Counter;

Task_Count : Task_Counter;

task type Traveler_Task is
   entry Start(T : in out Traveler);
end Traveler_Task;

task body Traveler_Task is
   T : Traveler;
begin
   accept Start(T : in out Traveler) do
      Task_Count.Increment;
      loop
         Move_Traveler(T);
         delay 1.0;
      end loop;
      Task_Count.Decrement;
   end Start;
end Traveler_Task;

task type Photo_Task is
   entry Start;
end Photo_Task;

task body Photo_Task is
begin
   accept Start do
      Task_Count.Increment;
      loop
         Take_Photo;
         delay 1.0;
      end loop;
      Task_Count.Decrement;
   end Start;
end Photo_Task;

Traveler_Tasks : array (1 .. 5) of Traveler_Task;
Photo_T : Photo_Task;

begin
   Reset(Gen);
   for I in Travelers'Range loop
      Travelers(I).ID := I;
      Travelers(I).Pos := (X => Integer(Random(Gen) * 10.0) + 1, Y => Integer(Random(Gen) * 10.0) + 1);
      Add_Traveler(Travelers(I));
      Traveler_Tasks(I).Start(Travelers(I));
   end loop;

   Photo_T.Start;

   -- Wait for all tasks to finish
   Task_Count.Wait_Until_Zero;
end Main;